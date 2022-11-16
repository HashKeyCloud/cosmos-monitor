package monitor

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"cosmosmonitor/config"
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/notification"
	"cosmosmonitor/rpc"
	"cosmosmonitor/types"
	"cosmosmonitor/utils"

	apolloDb "cosmosmonitor/db/apollo-db"
	bandDb "cosmosmonitor/db/band-db"
	cosmosDb "cosmosmonitor/db/cosmos-db"
	evmosDb "cosmosmonitor/db/evmos-db"
	injectiveDb "cosmosmonitor/db/injective-db"
	junoDb "cosmosmonitor/db/juno-db"
	neutronDb "cosmosmonitor/db/neutron-db"
	nyxDb "cosmosmonitor/db/nyx-db"
	persistenceDb "cosmosmonitor/db/persistence-db"
	providerDb "cosmosmonitor/db/provider-db"
	rizonDb "cosmosmonitor/db/rizon-db"
	secretDb "cosmosmonitor/db/secret-db"
	sommelierDb "cosmosmonitor/db/sommelier-db"
	sputnikDb "cosmosmonitor/db/sputnik-db"
	teritoriDb "cosmosmonitor/db/teritori-db"
	xplaDb "cosmosmonitor/db/xpla-db"

	apolloRpc "cosmosmonitor/rpc/apollo-rpc"
	bandRpc "cosmosmonitor/rpc/band-rpc"
	cosmosRpc "cosmosmonitor/rpc/cosmos-rpc"
	evmosRpc "cosmosmonitor/rpc/evmos-rpc"
	injectiveRpc "cosmosmonitor/rpc/injective-rpc"
	junoRpc "cosmosmonitor/rpc/juno-rpc"
	neutronRpc "cosmosmonitor/rpc/neutron-rpc"
	nyxRpc "cosmosmonitor/rpc/nyx-rpc"
	persistenceRpc "cosmosmonitor/rpc/persistence-rpc"
	providerRpc "cosmosmonitor/rpc/provider-rpc"
	rizonRpc "cosmosmonitor/rpc/rizon-rpc"
	secretRpc "cosmosmonitor/rpc/secret-rpc"
	sommelierRpc "cosmosmonitor/rpc/sommelier-rpc"
	sputnikRpc "cosmosmonitor/rpc/sputnik-rpc"
	teritoriRpc "cosmosmonitor/rpc/teritori-rpc"
	xplaRpc "cosmosmonitor/rpc/xpla-rpc"
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
	valRankingChan  chan []*types.ValRanking
}

var (
	preValJailed   = make(map[string]struct{}, 0)
	preValinActive = make(map[string]struct{}, 0)
	preProposalId  = make(map[string]struct{}, 0)
	monitorHeight  = make(map[string]int64, 0)
	startHeight    = make(map[string]int64, 0)
)

func NewMonitor() (*Monitor, error) {
	rpcClis := make(map[string]rpc.Client, 0)
	dbClis := make(map[string]db.DBCli, 0)
	apolloRpcCli, err := apolloRpc.InitApolloRpcCli()
	if err != nil {
		logger.Error("connect apollo rpc client error: ", err)
		return nil, err
	}
	bandRpcCli, err := bandRpc.InitBandRpcCli()
	if err != nil {
		logger.Error("connect band rpc client error: ", err)
		return nil, err
	}
	cosmosRpcCli, err := cosmosRpc.InitCosmosRpcCli()
	if err != nil {
		logger.Error("connect cosmos rpc client error: ", err)
		return nil, err
	}
	evmosRpcCli, err := evmosRpc.InitEvmosRpcCli()
	if err != nil {
		logger.Error("connect evmos rpc client error: ", err)
		return nil, err
	}

	injectiveRpcCli, err := injectiveRpc.InitInjectiveRpcCli()
	if err != nil {
		logger.Error("connect injective rpc client error: ", err)
		return nil, err
	}
	junoRpcCli, err := junoRpc.InitJunoRpcCli()
	if err != nil {
		logger.Error("connect juno rpc client error: ", err)
		return nil, err
	}
	neutronRpcCli, err := neutronRpc.InitNeutronRpcCli()
	if err != nil {
		logger.Error("connect neutron rpc client error: ", err)
		return nil, err
	}
	nyxRpcCli, err := nyxRpc.InitNyxRpcCli()
	if err != nil {
		logger.Error("connect nyx rpc client error: ", err)
		return nil, err
	}
	persistenceRpcCli, err := persistenceRpc.InitPersistenceRpcCli()
	if err != nil {
		logger.Error("connect persistence client error: ", err)
		return nil, err
	}
	providerRpcCli, err := providerRpc.InitProviderRpcCli()
	if err != nil {
		logger.Error("connect provider rpc client error: ", err)
		return nil, err
	}
	rizonRpcCli, err := rizonRpc.InitCosmosRpcCli()
	if err != nil {
		logger.Error("connect rizon rpc client error: ", err)
		return nil, err
	}

	secretRpcCli, err := secretRpc.InitCosmosRpcCli()
	if err != nil {
		logger.Error("connect secret rpc client error: ", err)
		return nil, err
	}
	sommelierRpcCli, err := sommelierRpc.InitSommelierRpcCli()
	if err != nil {
		logger.Error("connect sommelier rpc client error: ", err)
		return nil, err
	}
	sputnikRpcCli, err := sputnikRpc.InitSputnikRpcCli()
	if err != nil {
		logger.Error("connect sputnik rpc client error: ", err)
		return nil, err
	}
	teritoriRpcCli, err := teritoriRpc.InitTeritoriRpcCli()
	if err != nil {
		logger.Error("connect teritori rpc client error: ", err)
		return nil, err
	}
	xplaRpcCli, err := xplaRpc.InitXplaRpcCli()
	if err != nil {
		logger.Error("connect xpla rpc client error: ", err)
		return nil, err
	}

	rpcClis["apollo"] = apolloRpcCli
	rpcClis["cosmos"] = cosmosRpcCli
	rpcClis["band"] = bandRpcCli
	rpcClis["evmos"] = evmosRpcCli
	rpcClis["injective"] = injectiveRpcCli
	rpcClis["juno"] = junoRpcCli
	rpcClis["neutron"] = neutronRpcCli
	rpcClis["nyx"] = nyxRpcCli
	rpcClis["persistence"] = persistenceRpcCli
	rpcClis["provider"] = providerRpcCli
	rpcClis["rizon"] = rizonRpcCli
	rpcClis["secret"] = secretRpcCli
	rpcClis["sommelier"] = sommelierRpcCli
	rpcClis["sputnik"] = sputnikRpcCli
	rpcClis["teritori"] = teritoriRpcCli
	rpcClis["xpla"] = xplaRpcCli
	// init DB client
	apolloDbCli, err := apolloDb.InitApolloDbCli()
	if err != nil {
		logger.Error("connect apollo db client error:", err)
		return nil, err
	}
	bandDbCli, err := bandDb.InitBandDbCli()
	if err != nil {
		logger.Error("connect band db client error:", err)
		return nil, err
	}
	cosmosDbCli, err := cosmosDb.InitCosmosDbCli()
	if err != nil {
		logger.Error("connect cosmos db client error:", err)
		return nil, err
	}
	evmosDbCli, err := evmosDb.InitEvmosDbCli()
	if err != nil {
		logger.Error("connect evmos db client error:", err)
		return nil, err
	}
	injectiveDbCli, err := injectiveDb.InitInjectiveDbCli()
	if err != nil {
		logger.Error("connect injective db client error:", err)
		return nil, err
	}
	junoDbCli, err := junoDb.InitJunoDbCli()
	if err != nil {
		logger.Error("connect juno db client error:", err)
		return nil, err
	}
	neutronDbCli, err := neutronDb.InitNeutronDbCli()
	if err != nil {
		logger.Error("connect neutron db client error:", err)
		return nil, err
	}
	nyxDbCli, err := nyxDb.InitNyxDbCli()
	if err != nil {
		logger.Error("connect nyx db client error:", err)
		return nil, err
	}
	persistenceDbCli, err := persistenceDb.InitPersistenceDbCli()
	if err != nil {
		logger.Error("connect persistence db client error:", err)
		return nil, err
	}
	providerDbCli, err := providerDb.InitProviderDbCli()
	if err != nil {
		logger.Error("connect provider db client error:", err)
		return nil, err
	}
	rizonDbCli, err := rizonDb.InitRizonDbCli()
	if err != nil {
		logger.Error("connect rizon db client error:", err)
		return nil, err
	}
	secretDbCli, err := secretDb.InitSecretDbCli()
	if err != nil {
		logger.Error("connect secret db client error:", err)
		return nil, err
	}
	sommelierDbCli, err := sommelierDb.InitSommelierDbCli()
	if err != nil {
		logger.Error("connect sommelier db client error:", err)
		return nil, err
	}
	sputnikDbCli, err := sputnikDb.InitSputnikDbCli()
	if err != nil {
		logger.Error("connect sputnik db client error:", err)
		return nil, err
	}
	teritoriDbCli, err := teritoriDb.InitTeritoriDbCli()
	if err != nil {
		logger.Error("connect teritori db client error:", err)
		return nil, err
	}
	xplaDbCli, err := xplaDb.InitXplaDbCli()
	if err != nil {
		logger.Error("connect xpla db client error:", err)
		return nil, err
	}

	dbClis["apollo"] = apolloDbCli
	dbClis["band"] = bandDbCli
	dbClis["cosmos"] = cosmosDbCli
	dbClis["evmos"] = evmosDbCli
	dbClis["injective"] = injectiveDbCli
	dbClis["juno"] = junoDbCli
	dbClis["neutron"] = neutronDbCli
	dbClis["nyx"] = nyxDbCli
	dbClis["persistence"] = persistenceDbCli
	dbClis["provider"] = providerDbCli
	dbClis["rizon"] = rizonDbCli
	dbClis["secret"] = secretDbCli
	dbClis["sommelier"] = sommelierDbCli
	dbClis["sputnik"] = sputnikDbCli
	dbClis["teritori"] = teritoriDbCli
	dbClis["xpla"] = xplaDbCli

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
		valRankingChan:  make(chan []*types.ValRanking),
	}, nil
}

func (m *Monitor) Start() {
	epochTicker := time.NewTicker(time.Duration(viper.GetInt("alert.timeInterval")) * time.Second)
	consumerChains := make(map[string]string, 0)
	consumerChains["apollo"] = "provider"
	consumerChains["sputnik"] = "provider"
	for project := range m.RpcClis {
		startMonitorHeight := m.RpcClis[project].GetBlockHeight()
		startHeight[project] = startMonitorHeight
	}
	for range epochTicker.C {
		for project := range m.RpcClis {
			if project == "apollo" || project == "sputnik" {
				go m.consumer(project, consumerChains)
			} else {
				go m.provider(project)
			}
		}

	}
}

func (m *Monitor) consumer(consumerChain string, providerChain map[string]string) {
	mailSender := viper.GetString("mail.sender")
	receiver1Conf := fmt.Sprintf("mail.%sReceiver1", consumerChain)
	receiver2Conf := fmt.Sprintf("mail.%sReceiver2", consumerChain)
	receiver1 := viper.GetString(receiver1Conf)
	receiver2 := viper.GetString(receiver2Conf)
	mailReceiver := strings.Join([]string{receiver1, receiver2}, ",")

	mosDb, err := m.DbClis[providerChain[consumerChain]].GetMonitorObj()
	if len(mosDb) == 0 {
		time.Sleep(10 * time.Second)
		mosDb, err = m.DbClis[providerChain[consumerChain]].GetMonitorObj()
	}
	if err != nil {
		logger.Error("Failed to obtain production chain monitoring object, err:", err)
	}

	// list validator indices from config file
	logger.Info("Getting validators from config file.")
	operatorAddrs := config.GetOperatorAddrs(consumerChain)
	logger.Info("Getting Validators info from chain")

	mos := make(map[string]*types.MonitorObj)
	for _, moDb := range mosDb {
		mos[moDb.OperatorAddr] = moDb
	}

	mo := make([]*types.MonitorObj, 0)
	for _, operatorAddr := range operatorAddrs {
		if _, ok := mos[operatorAddr]; ok {
			mo = append(mo, &types.MonitorObj{
				Moniker:         mos[operatorAddr].Moniker,
				OperatorAddr:    mos[operatorAddr].OperatorAddr,
				OperatorAddrHex: mos[operatorAddr].OperatorAddrHex,
				SelfStakeAddr:   mos[operatorAddr].SelfStakeAddr,
			})
		} else {
			emailBody := fmt.Sprintf("If you want to monitor the validator %s of the consumer chain %s, you must monitor the consensus chain %s.",
				operatorAddr, consumerChain, providerChain[consumerChain])
			m.MailClient.SendMail(mailSender, mailReceiver, "Chain Exception", emailBody)
		}
	}
	careData := m.getValPerformanceRanking(consumerChain, mailSender, mailReceiver, mo)
	m.processData(&types.CaredData{
		ChainName:           consumerChain,
		ProposalAssignments: careData.ProposalAssignments,
		ValSigns:            careData.ValSigns,
		ValSignMisseds:      careData.ValSignMisseds,
		ValRankings:         careData.ValRankings,
	})

}

func (m *Monitor) provider(project string) {
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

	careData := m.getValPerformanceRanking(project, mailSender, mailReceiver, mo)

	m.processData(&types.CaredData{
		ChainName:           project,
		ValInfos:            valsInfo,
		Proposals:           proposals,
		ProposalAssignments: careData.ProposalAssignments,
		ValSigns:            careData.ValSigns,
		ValSignMisseds:      careData.ValSignMisseds,
		ValRankings:         careData.ValRankings,
	})
}

func (m *Monitor) getValPerformanceRanking(project, mailSender, mailReceiver string, mo []*types.MonitorObj) *types.CaredData {
	logger.Info("start getting validators performance")
	start, err := m.DbClis[project].GetBlockHeightFromDb()
	if err != nil {
		logger.Error("Failed to query block height from database，err:", err)
	}
	if start == 0 || startHeight[project]-start > 1000 {
		start = startHeight[project]
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
		}
	}
	logger.Info("Successfully get validators performance")

	logger.Info("start getting validators votingPower and ranking")
	valRanking, err := m.RpcClis[project].GetValRanking(mo, project)
	if err != nil {
		logger.Error("get proposal error: ", err)
		res := utils.Retry(func() bool {
			valRanking, err = m.RpcClis[project].GetValRanking(mo, project)
			if err != nil {
				return false
			} else {
				return true
			}
		}, []int{1, 3})
		if !res {
			emailBody := fmt.Sprintf("get cared data from %s RPC node error, please check.", project)
			m.MailClient.SendMail(mailSender, mailReceiver, "RPC Exception", emailBody)
		}
	}
	logger.Info("Successfully get validators votingPower and ranking")

	return &types.CaredData{
		ProposalAssignments: proposalAssignments,
		ValSigns:            valSign,
		ValSignMisseds:      valSignMissed,
		ValRankings:         valRanking,
	}
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

		end, err = m.DbClis[caredData.ChainName].GetBlockHeightFromDb()
		if err != nil {
			logger.Error("Failed to query block height from database，err:", err)
		}
		if end == 0 {
			end = startHeight[caredData.ChainName]
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
	if caredData.ValRankings != nil && len(caredData.ValRankings) > 0 {
		logger.Info("Start saving validator votingPower and ranking")
		err := m.DbClis[caredData.ChainName].BatchSaveRanking(caredData.ValRankings)
		if err != nil {
			logger.Error("save the validator votingPower and ranking fail:", err)
		}
		logger.Info("Save the validator votingPower and ranking successfully")

		valRankings := make([]*types.ValRanking, 0)

		projectRankingThreshold := fmt.Sprintf("alert.%sRankingThreshold", caredData.ChainName)
		rankingThreshold := viper.GetInt(projectRankingThreshold)
		for _, valRanking := range caredData.ValRankings {
			if valRanking.Ranking >= rankingThreshold {
				valRankings = append(valRankings, valRanking)
			}
		}
		m.valRankingChan <- valRankings
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

		select {
		case valRanking := <-m.valRankingChan:
			if len(valRanking) == 0 {
				break
			}
			va := notification.ParseValisRankingException(valRanking)
			err := m.MailClient.SendMail(mailSender, mailReceiver, va.Name(), va.Message())
			if err != nil {
				eventLogger.Error("send  validator ranking email error: ", err)
			}
			eventLogger.Info("send validator ranking email successful")
		case <-time.After(time.Second):
		}
	}
}

var logger = log.Logger.WithField("module", "monitor")
var eventLogger = log.EventLogger.WithField("module", "event")
