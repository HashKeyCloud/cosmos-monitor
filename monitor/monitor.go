package monitor

import (
	"cosmosmonitor/db"
	providerDb "cosmosmonitor/db/provider-db"
	injectiveRpc "cosmosmonitor/rpc/injective-rpc"
	provider "cosmosmonitor/rpc/provider-rpc"

	injectiveDb "cosmosmonitor/db/injective-db"
	junoDb "cosmosmonitor/db/juno-db"
	"cosmosmonitor/rpc"
	"cosmosmonitor/rpc/juno"
	"cosmosmonitor/utils"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"cosmosmonitor/config"
	cosmosDb "cosmosmonitor/db/cosmos-db"
	"cosmosmonitor/log"
	"cosmosmonitor/notification"
	cosmosRpc "cosmosmonitor/rpc/cosmos-rpc"
	"cosmosmonitor/types"
)

type Monitor struct {
	RpcClis         map[string]rpc.Client
	DbClis          map[string]db.DBCli
	MailClient      *notification.Client
	termChan        chan os.Signal
	valIsJailedChan chan []*types.ValIsJail
	valIsActiveChan chan []*types.ValIsActive
	missedSignChan  chan []*types.ValSignMissed
	proposalsChan   chan []*types.Proposal
}

var (
	preValJailed   = make(map[string]struct{}, 0)
	preValinActive = make(map[string]struct{}, 0)
	preProposalId  = make(map[string]struct{}, 0)
	monitorHeight  = make(map[string]int64, 0)
)

// 需要初始化不同的链RPC、数据库，邮件发送邮件可以一个
func NewMonitor() (*Monitor, error) {
	// init Cosmos RPC
	cosmosCli, err := cosmosRpc.InitCosmosRpcCli()
	if err != nil {
		logger.Error("connect cosmos rpc client error: ", err)
		return nil, err
	}
	injectiveCli, err := injectiveRpc.InitInjectiveRpcCli()
	if err != nil {
		logger.Error("connect cosmos rpc client error: ", err)
		return nil, err
	}

	junoCli, err := juno.InitJunoRpcCli()
	if err != nil {
		logger.Error("connect cosmos rpc client error: ", err)
		return nil, err
	}

	providerCli, err := provider.InitProviderRpcCli()
	if err != nil {
		logger.Error("connect cosmos rpc client error: ", err)
		return nil, err
	}

	rpcClis := make(map[string]rpc.Client)
	rpcClis["cosmos"] = cosmosCli
	rpcClis["injective"] = injectiveCli
	rpcClis["juno"] = junoCli
	rpcClis["provider"] = providerCli

	// init DB client
	cosmosDb, err := cosmosDb.InitCosmosDbCli()
	if err != nil {
		logger.Error("connect cosmos db client error:", err)
		return nil, err
	}

	injectiveDb, err := injectiveDb.InitInjectiveDbCli()
	if err != nil {
		logger.Error("connect injective db client error:", err)
		return nil, err
	}

	junoDb, err := junoDb.InitJunoDbCli()
	if err != nil {
		logger.Error("connect injective db client error:", err)
		return nil, err
	}

	providerDb, err := providerDb.InitProviderDbCli()
	if err != nil {
		logger.Error("connect injective db client error:", err)
		return nil, err
	}

	dbClis := make(map[string]db.DBCli)
	dbClis["cosmos"] = cosmosDb
	dbClis["injective"] = injectiveDb
	dbClis["juno"] = junoDb
	dbClis["providerDb"] = providerDb

	// init email client
	mailClient := notification.NewClient(
		viper.GetString("mail.host"),
		viper.GetInt("mail.port"),
		viper.GetString("mail.username"),
		viper.GetString("mail.Password"),
	)

	return &Monitor{
		RpcClis:         rpcClis,
		DbClis:          dbClis,
		MailClient:      mailClient,
		termChan:        make(chan os.Signal),
		valIsJailedChan: make(chan []*types.ValIsJail),
		missedSignChan:  make(chan []*types.ValSignMissed),
		proposalsChan:   make(chan []*types.Proposal),
		valIsActiveChan: make(chan []*types.ValIsActive),
	}, nil
}

func (m *Monitor) Start() {
	epochTicker := time.NewTicker(time.Duration(viper.GetInt("alert.timeInterval")) * time.Second)
	for range epochTicker.C {
		for project := range m.RpcClis {
			go m.start(project)
		}

	}
}

func (m *Monitor) start(project string) {
	mailSender := viper.GetString("mail.sender")
	receiver1Conf := fmt.Sprintf("mail.%sReceiver1", project)
	receiver2Conf := fmt.Sprintf("mail.%sReceiver2", project)
	receiver1 := viper.GetString(receiver1Conf)
	receiver2 := viper.GetString(receiver2Conf)
	mailReceiver := strings.Join([]string{receiver1, receiver2}, ",")

	// list validator indices from config file
	logger.Info("Getting validators from config file.")
	operatorAddrs := config.GetOperatorAddrs(project)
	logger.Info("Getting Validators info from chain")

	valsInfo, err := m.RpcClis[project].GetValInfo(operatorAddrs)
	if err != nil {
		res := utils.Retry(func() bool {
			valsInfo, err = m.RpcClis[project].GetValInfo(operatorAddrs)
			if err != nil {
				return false
			} else {
				return true
			}
		}, []int{1, 3})
		if !res {
			emailBody := fmt.Sprintf("get cared data from %s RPC node error, please check.", project)
			m.MailClient.SendMail(mailSender, mailReceiver, "RPC Exception", emailBody)
			return
		}
	}

	logger.Info("Successfully get validators!")

	mo := make([]*types.MonitorObj, 0)
	for _, valInfo := range valsInfo {
		mo = append(mo, &types.MonitorObj{
			Moniker:         valInfo.Moniker,
			OperatorAddr:    valInfo.OperatorAddr,
			OperatorAddrHex: valInfo.OperatorAddrHex,
			SelfStakeAddr:   valInfo.SelfStakeAddr,
		})
	}

	logger.Info("Start getting VOTING PERIOD proposals")
	proposals, err := m.RpcClis[project].GetProposal(mo)
	if err != nil {
		logger.Error("get proposal error: ", err)
		res := utils.Retry(func() bool {
			proposals, err = m.RpcClis[project].GetProposal(mo)
			if err != nil {
				return false
			} else {
				return true
			}
		}, []int{1, 3})
		if !res {
			emailBody := fmt.Sprintf("get cared data from %s RPC node error, please check.", project)
			m.MailClient.SendMail(mailSender, mailReceiver, "RPC Exception", emailBody)
			return
		}
	}
	logger.Info("Successfully get VOTING PERIOD proposals")

	logger.Info("start getting validators performance")
	start, err := m.DbClis[project].GetBlockHeightFromDb(project)
	if err != nil {
		logger.Error("Failed to query block height from database，err:", err)
	}
	proposalAssignments, valSign, valSignMissed, err := m.RpcClis[project].GetValPerformance(start, mo)
	if err != nil {
		logger.Error("get proposal error: ", err)
		res := utils.Retry(func() bool {
			proposalAssignments, valSign, valSignMissed, err = m.RpcClis[project].GetValPerformance(start, mo)
			if err != nil {
				return false
			} else {
				return true
			}
		}, []int{1, 3})
		if !res {
			emailBody := fmt.Sprintf("get cared data from %s RPC node error, please check.", project)
			m.MailClient.SendMail(mailSender, mailReceiver, "RPC Exception", emailBody)
			return
		}
	}
	logger.Info("Successfully get validators performance")

	m.processData(&types.CaredData{
		ChainName:           project,
		ValInfos:            valsInfo,
		Proposals:           proposals,
		ProposalAssignments: proposalAssignments,
		ValSigns:            valSign,
		ValSignMisseds:      valSignMissed,
	})
}

func (m *Monitor) WaitInterrupted() {
	<-m.termChan
	logger.Info("monitor shutdown signal received")
}

func (m *Monitor) processData(caredData *types.CaredData) {
	var end int64
	if caredData.ValInfos != nil && len(caredData.ValInfos) > 0 {
		logger.Infof("Start saving %s validator information\n", caredData.ChainName)
		err := m.DbClis[caredData.ChainName].SaveValInfo(caredData.ValInfos)
		if err != nil {
			logger.Errorf("save %s valdator info fail, err:%v\n", caredData.ChainName, err)
		}
		logger.Infof("Save the %s validator information successfully\n", caredData.ChainName)

		valIsJailed := make([]*types.ValIsJail, 0)
		valIsActive := make([]*types.ValIsActive, 0)
		newpreValJailed := make(map[string]struct{}, 0)
		newpreValinActive := make(map[string]struct{}, 0)

		for _, valInfo := range caredData.ValInfos {
			if _, ok := preValJailed[caredData.ChainName+valInfo.Moniker]; !ok && valInfo.Jailed {
				valIsJailed = append(valIsJailed, &types.ValIsJail{
					ChainName: caredData.ChainName,
					Moniker:   valInfo.Moniker,
				})
			}
			if _, ok := preValinActive[caredData.ChainName+valInfo.Moniker]; !ok && valInfo.Status != 3 {
				valIsActive = append(valIsActive, &types.ValIsActive{
					ChainName: caredData.ChainName,
					Moniker:   valInfo.Moniker,
					Status:    valInfo.Status,
				})
			}
			if valInfo.Jailed {
				newpreValJailed[caredData.ChainName+valInfo.Moniker] = struct{}{}
			}
			if valInfo.Status != 3 {
				newpreValinActive[caredData.ChainName+valInfo.Moniker] = struct{}{}
			}
		}
		preValJailed = newpreValJailed
		preValinActive = newpreValinActive

		m.valIsJailedChan <- valIsJailed
		m.valIsActiveChan <- valIsActive
	}

	if caredData.Proposals != nil && len(caredData.Proposals) > 0 {
		logger.Info("Start saving proposals information")
		err := m.DbClis[caredData.ChainName].BatchSaveProposals(caredData.Proposals)
		if err != nil {
			logger.Error("save the proposals information fail:", err)
		}
		logger.Info("Save the proposals information successfully")

		proposals := make([]*types.Proposal, 0)
		newPreProposalId := make(map[string]struct{}, 0)

		for _, proposal := range caredData.Proposals {
			proposalFlag := caredData.ChainName + strconv.Itoa(int(proposal.ProposalId))
			if _, ok := newPreProposalId[proposalFlag]; !ok {
				newPreProposalId[proposalFlag] = struct{}{}
			}

			if _, ok := preProposalId[proposalFlag]; !ok {
				proposals = append(proposals, &types.Proposal{
					ChainName:       caredData.ChainName,
					ProposalId:      proposal.ProposalId,
					VotingStartTime: proposal.VotingStartTime,
					VotingEndTime:   proposal.VotingEndTime,
					Description:     proposal.Description,
					Moniker:         proposal.Moniker,
					OperatorAddr:    proposal.OperatorAddr,
					Status:          proposal.Status,
				})
				preProposalId[proposalFlag] = struct{}{}
			}
		}
		preProposalId = newPreProposalId
		m.proposalsChan <- proposals
	}

	if caredData.ProposalAssignments != nil && len(caredData.ProposalAssignments) > 0 {
		logger.Info("Start saving proposal assignments information")
		err := m.DbClis[caredData.ChainName].BatchSaveProposalAssignments(caredData.ProposalAssignments)
		if err != nil {
			logger.Error("save proposal assignments fail:", err)
		}
		logger.Info("Save the proposal assignments successfully")
	}

	if caredData.ValSigns != nil && len(caredData.ValSigns) > 0 {
		logger.Info("Start saving validator signs")
		err := m.DbClis[caredData.ChainName].BatchSaveValSign(caredData.ValSigns)
		if err != nil {
			logger.Error("save validator sign fail:", err)
		}
		logger.Info("Save the validator sign successfully")
	}

	if caredData.ValSignMisseds != nil && len(caredData.ValSignMisseds) > 0 {
		logger.Info("Start saving validator sign misseds")
		err := m.DbClis[caredData.ChainName].BatchSaveValSignMissed(caredData.ValSignMisseds)
		if err != nil {
			logger.Error("save validator sign missed fail:", err)
		}
		logger.Info("Save the validator sign missed successfully")

		end, err = m.DbClis[caredData.ChainName].GetBlockHeightFromDb(caredData.ChainName)
		if err != nil {
			logger.Error("Failed to query block height from database，err:", err)
		}
		interval := viper.GetInt("alert.blockInterval")
		start := end - int64(interval) + int64(1)
		valSignMissed, err := m.DbClis[caredData.ChainName].GetValSignMissedFromDb(start, end)
		if err != nil {
			logger.Error("Failed to query validator sign missed from database, err:", err)
		}
		valSignMissedMap := make(map[string][]int)
		for _, signMissed := range valSignMissed {
			valSignMissedMap[signMissed.OperatorAddr] = append(valSignMissedMap[signMissed.OperatorAddr], signMissed.BlockHeight)
		}
		valsMoniker, err := m.DbClis[caredData.ChainName].GetValMoniker()
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
					ChainName:    caredData.ChainName,
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
							ChainName:    caredData.ChainName,
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
	if monitorHeight[caredData.ChainName] != endHeight {
		m.DbClis[caredData.ChainName].BatchSaveValStats(endHeight-int64(timeInterval)+int64(1), endHeight)
		monitorHeight[caredData.ChainName] = endHeight
	}
}

func (m *Monitor) SendEmail() {
	for project := range m.RpcClis {
		mailSender := viper.GetString("mail.sender")
		receiver1Conf := fmt.Sprintf("mail.%sReceiver1", project)
		receiver2Conf := fmt.Sprintf("mail.%sReceiver2", project)
		receiver1 := viper.GetString(receiver1Conf)
		receiver2 := viper.GetString(receiver2Conf)
		mailReceiver := strings.Join([]string{receiver1, receiver2}, ",")
		go m.sendEmail(mailSender, receiver1, mailReceiver)
	}
}

func (m *Monitor) sendEmail(mailSender, receiver1, mailReceiver string) {
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
