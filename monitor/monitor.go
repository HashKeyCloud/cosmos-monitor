package monitor

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"

	"cosmosmonitor/config"
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/notification"
	cosmosRpc "cosmosmonitor/rpc/cosmos"
	"cosmosmonitor/types"
	"cosmosmonitor/utils"
)

type Monitor struct {
	CosmosRpcCli    *cosmosRpc.CosmosCli
	CosmosDbCli     *db.DbCli
	MailClient      *notification.Client
	termChan        chan os.Signal
	valIsJailedChan chan []string
	missedSignChan  chan []*types.ValSignMissed
	proposalsChan   chan []*types.Proposal
	valIsActiveChan chan []*types.ValIsActive
}

var (
	preValJailed         = make(map[string]struct{}, 0)
	preValinActive       = make(map[string]struct{}, 0)
	preProposalId        = make(map[int64]struct{}, 0)
	monitorHeight  int64 = 0
)

func NewMonitor() (*Monitor, error) {
	// init Cosmos RPC
	rpcEndpoint := fmt.Sprintf("%s:%s", viper.GetString("cosmos.ip"), viper.GetString("cosmos.gRPCport"))
	cosmosCli, err := cosmosRpc.NewCosmosRpcCli(rpcEndpoint)
	if err != nil {
		logger.Error("connect rpc client error: ", err)
		return nil, err
	}
	// init DB client
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		Name:     viper.GetString("postgres.name"),
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect database server error: ", err)
	}
	// init email client
	mailClient := notification.NewClient(
		viper.GetString("mail.host"),
		viper.GetInt("mail.port"),
		viper.GetString("mail.username"),
		viper.GetString("mail.Password"),
	)

	return &Monitor{
		CosmosRpcCli:    cosmosCli,
		CosmosDbCli:     &db.DbCli{Conn: dbCli},
		MailClient:      mailClient,
		termChan:        make(chan os.Signal),
		valIsJailedChan: make(chan []string),
		missedSignChan:  make(chan []*types.ValSignMissed),
		proposalsChan:   make(chan []*types.Proposal),
		valIsActiveChan: make(chan []*types.ValIsActive),
	}, nil
}

func (m *Monitor) Start() {
	mailSender := viper.GetString("mail.sender")
	receiver1 := viper.GetString("mail.receiver1")
	receiver2 := viper.GetString("mail.receiver2")
	mailReceiver := strings.Join([]string{receiver1, receiver2}, ",")

	epochTicker := time.NewTicker(time.Duration(viper.GetInt("alert.timeInterval")) * time.Second)
	for range epochTicker.C {
		// list validator indices from config file
		logger.Info("Getting validators from config file.")
		operatorAdds := config.GetoperatorAddrs()

		logger.Info("Getting Validators info from chain")
		valsInfo, err := m.CosmosRpcCli.GetValInfo(operatorAdds)
		if err != nil {
			logger.Error("get cared data error: ", err)
			res := utils.Retry(func() bool {
				valsInfo, err = m.CosmosRpcCli.GetValInfo(operatorAdds)
				if err != nil {
					return false
				} else {
					return true
				}
			}, []int{1, 3})
			if !res {
				m.MailClient.SendMail(mailSender, mailReceiver, "RPC Exception", "get cared data from RPC node error, please check.")
				continue
			}
		}
		logger.Info("Successfully get validators!")

		mo, err := m.CosmosDbCli.GetMonitorObj()
		if err != nil {
			logger.Error("")
		}
		logger.Info("Start getting VOTING PERIOD proposals")
		proposals, err := m.CosmosRpcCli.GetProposal(mo)
		if err != nil {
			logger.Error("get proposal error: ", err)
			res := utils.Retry(func() bool {
				proposals, err = m.CosmosRpcCli.GetProposal(mo)
				if err != nil {
					return false
				} else {
					return true
				}
			}, []int{1, 3})
			if !res {
				m.MailClient.SendMail(mailSender, mailReceiver, "RPC Exception", "get cared data from RPC node error, please check.")
				continue
			}
		}
		logger.Info("Successfully get VOTING PERIOD proposals")

		logger.Info("start getting validators performance")
		start, err := m.CosmosDbCli.GetBlockHeightFromDb()
		if err != nil {
			logger.Error("Failed to query block height from database，err:", err)
		}
		proposalAssignments, valSign, valSignMissed, err := m.CosmosRpcCli.GetValPerformance(start, mo)
		if err != nil {
			logger.Error("get proposal error: ", err)
			res := utils.Retry(func() bool {
				proposalAssignments, valSign, valSignMissed, err = m.CosmosRpcCli.GetValPerformance(start, mo)
				if err != nil {
					return false
				} else {
					return true
				}
			}, []int{1, 3})
			if !res {
				m.MailClient.SendMail(mailSender, mailReceiver, "RPC Exception", "get cared data from RPC node error, please check.")
				continue
			}
		}
		logger.Info("Successfully get validators performance")

		m.processData(&types.CaredData{
			ValInfos:            valsInfo,
			Proposals:           proposals,
			ProposalAssignments: proposalAssignments,
			ValSigns:            valSign,
			ValSignMisseds:      valSignMissed,
		})

	}
}

func (m *Monitor) WaitInterrupted() {
	<-m.termChan
	logger.Info("monitor shutdown signal received")
}

func (m *Monitor) processData(caredData *types.CaredData) {
	var end int64
	if caredData.ValInfos != nil && len(caredData.ValInfos) > 0 {
		logger.Info("Start saving validator information")
		err := m.CosmosDbCli.SaveValInfo(caredData.ValInfos)
		if err != nil {
			logger.Error("save valdator info fail:", err)
		}
		logger.Info("Save the validator information successfully")

		valIsJailed := make([]string, 0)
		valIsActive := make([]*types.ValIsActive, 0)
		newpreValJailed := make(map[string]struct{}, 0)
		newpreValinActive := make(map[string]struct{}, 0)

		for _, valInfo := range caredData.ValInfos {
			if _, ok := preValJailed[valInfo.Moniker]; !ok && valInfo.Jailed {
				valIsJailed = append(valIsJailed, valInfo.Moniker)
			}
			if _, ok := preValinActive[valInfo.Moniker]; !ok && valInfo.Status != 3 {
				valIsActive = append(valIsActive, &types.ValIsActive{
					Moniker: valInfo.Moniker,
					Status:  valInfo.Status,
				})
			}
			if valInfo.Jailed {
				newpreValJailed[valInfo.Moniker] = struct{}{}
			}
			if valInfo.Status != 3 {
				newpreValinActive[valInfo.Moniker] = struct{}{}
			}
		}
		preValJailed = newpreValJailed
		preValinActive = newpreValinActive

		m.valIsJailedChan <- valIsJailed
		m.valIsActiveChan <- valIsActive
	}

	if caredData.Proposals != nil && len(caredData.Proposals) > 0 {
		logger.Info("Start saving proposals information")
		err := m.CosmosDbCli.BatchSaveProposals(caredData.Proposals)
		if err != nil {
			logger.Error("save the proposals information fail:", err)
		}
		logger.Info("Save the proposals information successfully")

		proposals := make([]*types.Proposal, 0)
		newPreProposalId := make(map[int64]struct{}, 0)

		for _, proposal := range caredData.Proposals {
			if _, ok := newPreProposalId[proposal.ProposalId]; !ok {
				newPreProposalId[proposal.ProposalId] = struct{}{}
			}

			if _, ok := preProposalId[proposal.ProposalId]; !ok {
				proposals = append(proposals, proposal)
				preProposalId[proposal.ProposalId] = struct{}{}
			}
		}
		preProposalId = newPreProposalId
		m.proposalsChan <- proposals
	}

	if caredData.ProposalAssignments != nil && len(caredData.ProposalAssignments) > 0 {
		logger.Info("Start saving proposal assignments information")
		err := m.CosmosDbCli.BatchSaveProposalAssignments(caredData.ProposalAssignments)
		if err != nil {
			logger.Error("save proposal assignments fail:", err)
		}
		logger.Info("Save the proposal assignments successfully")
	}

	if caredData.ValSigns != nil && len(caredData.ValSigns) > 0 {
		logger.Info("Start saving validator signs")
		err := m.CosmosDbCli.BatchSaveValSign(caredData.ValSigns)
		if err != nil {
			logger.Error("save validator sign fail:", err)
		}
		logger.Info("Save the validator sign successfully")
	}

	if caredData.ValSignMisseds != nil && len(caredData.ValSignMisseds) > 0 {
		logger.Info("Start saving validator sign misseds")
		err := m.CosmosDbCli.BatchSaveValSignMissed(caredData.ValSignMisseds)
		if err != nil {
			logger.Error("save validator sign missed fail:", err)
		}
		logger.Info("Save the validator sign missed successfully")

		end, err = m.CosmosDbCli.GetBlockHeightFromDb()
		if err != nil {
			logger.Error("Failed to query block height from database，err:", err)
		}
		interval := viper.GetInt("alert.blockInterval")
		start := end - int64(interval) + int64(1)
		valSignMissed, err := m.CosmosDbCli.GetValSignMissedFromDb(start, end)
		if err != nil {
			logger.Error("Failed to query validator sign missed from database, err:", err)
		}
		valSignMissedMap := make(map[string][]int)
		for _, signMissed := range valSignMissed {
			valSignMissedMap[signMissed.OperatorAddr] = append(valSignMissedMap[signMissed.OperatorAddr], signMissed.BlockHeight)
		}
		valsMoniker, err := m.CosmosDbCli.GetValMoniker()
		if err != nil {
			logger.Error("Failed to query validator moniker from database, err:", err)
		}
		valMonikerMap := make(map[string]string)
		for _, valMoniker := range valsMoniker {
			valMonikerMap[valMoniker.OperatorAddr] = valMoniker.Moniker
		}
		proportion := viper.GetFloat64("alert.proportion")
		missedSign := make([]*types.ValSignMissed, 0)
		recorded := make(map[string]struct{}, 0)

		for operatorAddr, missedBlcoks := range valSignMissedMap {
			if float64(len(missedBlcoks))/float64(interval) > proportion {
				missedSign = append(missedSign, &types.ValSignMissed{
					OperatorAddr: valMonikerMap[operatorAddr],
					BlockHeight:  int(end),
				})
				recorded[operatorAddr] = struct{}{}
			}

			if len(missedBlcoks) > 5 {
				sort.Ints(missedBlcoks)
				for i := 0; i < len(missedBlcoks)-5; i++ {
					if _, ok := recorded[operatorAddr]; !ok && missedBlcoks[i+4]-missedBlcoks[i] == 4 {
						missedSign = append(missedSign, &types.ValSignMissed{
							OperatorAddr: valMonikerMap[operatorAddr],
							BlockHeight:  int(end),
						})
						recorded[operatorAddr] = struct{}{}
					}
				}
			}

		}
		m.missedSignChan <- missedSign
	}

	timeInterval := viper.GetInt("alert.timeInterval")
	endHeight := end / int64(timeInterval) * int64(timeInterval)
	if monitorHeight != endHeight {
		m.CosmosDbCli.BatchSaveValStats(endHeight-int64(timeInterval)+int64(1), endHeight)
		monitorHeight = endHeight
	}
}

func (m *Monitor) ProcessData() {
	mailSender := viper.GetString("mail.sender")
	receiver1 := viper.GetString("mail.receiver1")
	receiver2 := viper.GetString("mail.receiver2")
	mailReceiver := strings.Join([]string{receiver1, receiver2}, ",")
	for {
		select {
		case valJailed := <-m.valIsJailedChan:
			if len(valJailed) == 0 {
				break
			}
			vj := notification.ParseValJailedException(valJailed)
			err := m.MailClient.SendMail(mailSender, mailReceiver, vj.Name(), vj.Message())
			if err != nil {
				eventLogger.Error("send validator jailed email error: ", err)
			}
			eventLogger.Info("send validator jailed email successful")
		case <-time.After(time.Second):
		}
		select {
		case valisAtive := <-m.valIsActiveChan:
			if len(valisAtive) == 0 {
				break
			}
			va := notification.ParseValisActiveException(valisAtive)
			err := m.MailClient.SendMail(mailSender, mailReceiver, va.Name(), va.Message())
			if err != nil {
				eventLogger.Error("send  validator inActive email error: ", err)
			}
			eventLogger.Info("send validator inActive email successful")
		case <-time.After(time.Second):
		}
		select {
		case missedSign := <-m.missedSignChan:
			if len(missedSign) == 0 {
				break
			}
			va := notification.ParseSyncException(missedSign)
			err := m.MailClient.SendMail(mailSender, mailReceiver, va.Name(), va.Message())
			if err != nil {
				eventLogger.Error("send sign missed email error: ", err)
			}
			eventLogger.Info("send validator sign missed email successful")
		case <-time.After(time.Second):
		}
		select {
		case proposals := <-m.proposalsChan:
			if len(proposals) == 0 {
				break
			}
			p := notification.ParseProposalException(proposals)
			err := m.MailClient.SendMail(mailSender, receiver1, p.Name(), p.Message())
			if err != nil {
				eventLogger.Error("send proposal email error: ", err)
			}
			eventLogger.Info("send proposal email successful")
		case <-time.After(time.Second):
		}
	}
}

var logger = log.Logger.WithField("module", "monitor")
var eventLogger = log.EventLogger.WithField("module", "event")
