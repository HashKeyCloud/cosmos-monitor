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
	"cosmosmonitor/rpc"
	"cosmosmonitor/types"
	"cosmosmonitor/utils"

	acrechainDb "cosmosmonitor/db/acrechain-db"
	akashDb "cosmosmonitor/db/akash-db"
	apolloDb "cosmosmonitor/db/apollo-db"
	axelarDb "cosmosmonitor/db/axelar-db"
	bandDb "cosmosmonitor/db/band-db"
	cosmosDb "cosmosmonitor/db/cosmos-db"
	evmosDb "cosmosmonitor/db/evmos-db"
	gopherDb "cosmosmonitor/db/gopher-db"
	heroDb "cosmosmonitor/db/hero-db"
	injectiveDb "cosmosmonitor/db/injective-db"
	junoDb "cosmosmonitor/db/juno-db"
	neutron_consumerDb "cosmosmonitor/db/neutron-consumer-db"
	neutronDb "cosmosmonitor/db/neutron-db"
	nyxDb "cosmosmonitor/db/nyx-db"
	okp4Db "cosmosmonitor/db/okp4-db"
	persistenceDb "cosmosmonitor/db/persistence-db"
	providerDb "cosmosmonitor/db/provider-db"
	rizonDb "cosmosmonitor/db/rizon-db"
	secretDb "cosmosmonitor/db/secret-db"
	sommelierDb "cosmosmonitor/db/sommelier-db"
	sputnikDb "cosmosmonitor/db/sputnik-db"
	teritoriDb "cosmosmonitor/db/teritori-db"
	xplaDb "cosmosmonitor/db/xpla-db"
	gopherRpc "cosmosmonitor/rpc/gopher-rpc"
	heroRpc "cosmosmonitor/rpc/hero-rpc"

	acrechainRpc "cosmosmonitor/rpc/acrechain-rpc"
	akashRpc "cosmosmonitor/rpc/akash-rpc"
	apolloRpc "cosmosmonitor/rpc/apollo-rpc"
	axelarRpc "cosmosmonitor/rpc/axelar-rpc"
	bandRpc "cosmosmonitor/rpc/band-rpc"
	cosmosRpc "cosmosmonitor/rpc/cosmos-rpc"
	evmosRpc "cosmosmonitor/rpc/evmos-rpc"
	injectiveRpc "cosmosmonitor/rpc/injective-rpc"
	junoRpc "cosmosmonitor/rpc/juno-rpc"
	neutron_consumerRpc "cosmosmonitor/rpc/neutron-consumer-rpc"
	neutronRpc "cosmosmonitor/rpc/neutron-rpc"
	nyxRpc "cosmosmonitor/rpc/nyx-rpc"
	okp4Rpc "cosmosmonitor/rpc/okp4-rpc"
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
	valIsJailedChan map[string]chan []*types.ValIsJail
	valIsActiveChan map[string]chan []*types.ValIsActive
	missedSignChan  map[string]chan []*types.ValSignMissed
	proposalsChan   map[string]chan []*types.Proposal
	valRankingChan  map[string]chan []*types.ValRanking
}

var (
	preValJailed   = make(map[string]map[string]struct{}, 0)
	preValinActive = make(map[string]map[string]struct{}, 0)
	preProposalId  = make(map[string]map[int64]struct{}, 0)
	monitorHeight  = make(map[string]int64, 0)
	startHeight    = make(map[string]int64, 0)
)

func NewMonitor() (*Monitor, error) {
	rpcClis := make(map[string]rpc.Client, 0)
	dbClis := make(map[string]db.DBCli, 0)

	valIsJailedChan := make(map[string]chan []*types.ValIsJail)
	missedSignChan := make(map[string]chan []*types.ValSignMissed)
	proposalsChan := make(map[string]chan []*types.Proposal)
	valIsActiveChan := make(map[string]chan []*types.ValIsActive)
	valRankingChan := make(map[string]chan []*types.ValRanking)
	if viper.GetBool("alert.acrechainIsMonitored") {
		acrechainRpcCli, err := acrechainRpc.InitAcrechainRpcCli()
		if err != nil {
			logger.Error("connect acrechain rpc client error: ", err)
			return nil, err
		}
		rpcClis["acrechain"] = acrechainRpcCli
		acrechainDbCli, err := acrechainDb.InitAcrechainDbCli()
		if err != nil {
			logger.Error("connect acrechain db client error:", err)
			return nil, err
		}
		dbClis["acrechain"] = acrechainDbCli
		valIsJailedChan["acrechain"] = make(chan []*types.ValIsJail)
		missedSignChan["acrechain"] = make(chan []*types.ValSignMissed)
		proposalsChan["acrechain"] = make(chan []*types.Proposal)
		valIsActiveChan["acrechain"] = make(chan []*types.ValIsActive)
		valRankingChan["acrechain"] = make(chan []*types.ValRanking)
		preValJailed["acrechain"] = make(map[string]struct{}, 0)
		preValinActive["acrechain"] = make(map[string]struct{}, 0)
		preProposalId["acrechain"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.akashIsMonitored") {
		akashRpcCli, err := akashRpc.InitAkashRpcCli()
		if err != nil {
			logger.Error("connect akash rpc client error: ", err)
			return nil, err
		}
		rpcClis["akash"] = akashRpcCli
		akashDbCli, err := akashDb.InitAkashDbCli()
		if err != nil {
			logger.Error("connect akash db client error:", err)
			return nil, err
		}
		dbClis["akash"] = akashDbCli
		valIsJailedChan["akash"] = make(chan []*types.ValIsJail)
		missedSignChan["akash"] = make(chan []*types.ValSignMissed)
		proposalsChan["akash"] = make(chan []*types.Proposal)
		valIsActiveChan["akash"] = make(chan []*types.ValIsActive)
		valRankingChan["akash"] = make(chan []*types.ValRanking)
		preValJailed["akash"] = make(map[string]struct{}, 0)
		preValinActive["akash"] = make(map[string]struct{}, 0)
		preProposalId["akash"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.apolloIsMonitored") {
		apolloRpcCli, err := apolloRpc.InitApolloRpcCli()
		if err != nil {
			logger.Error("connect apollo rpc client error: ", err)
			return nil, err
		}
		rpcClis["apollo"] = apolloRpcCli
		apolloDbCli, err := apolloDb.InitApolloDbCli()
		if err != nil {
			logger.Error("connect apollo db client error:", err)
			return nil, err
		}
		dbClis["apollo"] = apolloDbCli
		valIsJailedChan["apollo"] = make(chan []*types.ValIsJail)
		missedSignChan["apollo"] = make(chan []*types.ValSignMissed)
		proposalsChan["apollo"] = make(chan []*types.Proposal)
		valIsActiveChan["apollo"] = make(chan []*types.ValIsActive)
		valRankingChan["apollo"] = make(chan []*types.ValRanking)
		preValJailed["apollo"] = make(map[string]struct{}, 0)
		preValinActive["apollo"] = make(map[string]struct{}, 0)
		preProposalId["apollo"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.axelarIsMonitored") {
		axelarRpcCli, err := axelarRpc.InitAxelarRpcCli()
		if err != nil {
			logger.Error("connect axelar rpc client error: ", err)
			return nil, err
		}
		rpcClis["axelar"] = axelarRpcCli
		axelarDbCli, err := axelarDb.InitAxelarDbCli()
		if err != nil {
			logger.Error("connect axelar db client error:", err)
			return nil, err
		}
		dbClis["axelar"] = axelarDbCli
		valIsJailedChan["axelar"] = make(chan []*types.ValIsJail)
		missedSignChan["axelar"] = make(chan []*types.ValSignMissed)
		proposalsChan["axelar"] = make(chan []*types.Proposal)
		valIsActiveChan["axelar"] = make(chan []*types.ValIsActive)
		valRankingChan["axelar"] = make(chan []*types.ValRanking)
		preValJailed["axelar"] = make(map[string]struct{}, 0)
		preValinActive["axelar"] = make(map[string]struct{}, 0)
		preProposalId["axelar"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.bandIsMonitored") {
		bandRpcCli, err := bandRpc.InitBandRpcCli()
		if err != nil {
			logger.Error("connect band rpc client error: ", err)
			return nil, err
		}
		rpcClis["band"] = bandRpcCli
		bandDbCli, err := bandDb.InitBandDbCli()
		if err != nil {
			logger.Error("connect band db client error:", err)
			return nil, err
		}
		dbClis["band"] = bandDbCli
		valIsJailedChan["band"] = make(chan []*types.ValIsJail)
		missedSignChan["band"] = make(chan []*types.ValSignMissed)
		proposalsChan["band"] = make(chan []*types.Proposal)
		valIsActiveChan["band"] = make(chan []*types.ValIsActive)
		valRankingChan["band"] = make(chan []*types.ValRanking)
		preValJailed["band"] = make(map[string]struct{}, 0)
		preValinActive["band"] = make(map[string]struct{}, 0)
		preProposalId["band"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.cosmosIsMonitored") {
		cosmosRpcCli, err := cosmosRpc.InitCosmosRpcCli()
		if err != nil {
			logger.Error("connect cosmos rpc client error: ", err)
			return nil, err
		}
		rpcClis["cosmos"] = cosmosRpcCli
		cosmosDbCli, err := cosmosDb.InitCosmosDbCli()
		if err != nil {
			logger.Error("connect cosmos db client error:", err)
			return nil, err
		}
		dbClis["cosmos"] = cosmosDbCli
		valIsJailedChan["cosmos"] = make(chan []*types.ValIsJail)
		missedSignChan["cosmos"] = make(chan []*types.ValSignMissed)
		proposalsChan["cosmos"] = make(chan []*types.Proposal)
		valIsActiveChan["cosmos"] = make(chan []*types.ValIsActive)
		valRankingChan["cosmos"] = make(chan []*types.ValRanking)
		preValJailed["cosmos"] = make(map[string]struct{}, 0)
		preValinActive["cosmos"] = make(map[string]struct{}, 0)
		preProposalId["cosmos"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.evmosIsMonitored") {
		evmosRpcCli, err := evmosRpc.InitEvmosRpcCli()
		if err != nil {
			logger.Error("connect evmos rpc client error: ", err)
			return nil, err
		}
		rpcClis["evmos"] = evmosRpcCli
		evmosDbCli, err := evmosDb.InitEvmosDbCli()
		if err != nil {
			logger.Error("connect evmos db client error:", err)
			return nil, err
		}
		dbClis["evmos"] = evmosDbCli
		valIsJailedChan["evmos"] = make(chan []*types.ValIsJail)
		missedSignChan["evmos"] = make(chan []*types.ValSignMissed)
		proposalsChan["evmos"] = make(chan []*types.Proposal)
		valIsActiveChan["evmos"] = make(chan []*types.ValIsActive)
		valRankingChan["evmos"] = make(chan []*types.ValRanking)
		preValJailed["evmos"] = make(map[string]struct{}, 0)
		preValinActive["evmos"] = make(map[string]struct{}, 0)
		preProposalId["evmos"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.gopherIsMonitored") {
		gopherRpcCli, err := gopherRpc.InitGopherRpcCli()
		if err != nil {
			logger.Error("connect gopher rpc client error: ", err)
			return nil, err
		}
		rpcClis["gopher"] = gopherRpcCli
		gopherDbCli, err := gopherDb.InitGopherDbCli()
		if err != nil {
			logger.Error("connect gopher db client error:", err)
			return nil, err
		}
		dbClis["gopher"] = gopherDbCli
		valIsJailedChan["gopher"] = make(chan []*types.ValIsJail)
		missedSignChan["gopher"] = make(chan []*types.ValSignMissed)
		proposalsChan["gopher"] = make(chan []*types.Proposal)
		valIsActiveChan["gopher"] = make(chan []*types.ValIsActive)
		valRankingChan["gopher"] = make(chan []*types.ValRanking)
		preValJailed["gopher"] = make(map[string]struct{}, 0)
		preValinActive["gopher"] = make(map[string]struct{}, 0)
		preProposalId["gopher"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.heroIsMonitored") {
		heroRpcCli, err := heroRpc.InitHeroRpcCli()
		if err != nil {
			logger.Error("connect hero rpc client error: ", err)
			return nil, err
		}
		rpcClis["hero"] = heroRpcCli
		heroDbCli, err := heroDb.InitHeroDbCli()
		if err != nil {
			logger.Error("connect hero db client error:", err)
			return nil, err
		}
		dbClis["hero"] = heroDbCli
		valIsJailedChan["hero"] = make(chan []*types.ValIsJail)
		missedSignChan["hero"] = make(chan []*types.ValSignMissed)
		proposalsChan["hero"] = make(chan []*types.Proposal)
		valIsActiveChan["hero"] = make(chan []*types.ValIsActive)
		valRankingChan["hero"] = make(chan []*types.ValRanking)
		preValJailed["hero"] = make(map[string]struct{}, 0)
		preValinActive["hero"] = make(map[string]struct{}, 0)
		preProposalId["hero"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.injectiveIsMonitored") {
		injectiveRpcCli, err := injectiveRpc.InitInjectiveRpcCli()
		if err != nil {
			logger.Error("connect injective rpc client error: ", err)
			return nil, err
		}
		rpcClis["injective"] = injectiveRpcCli
		injectiveDbCli, err := injectiveDb.InitInjectiveDbCli()
		if err != nil {
			logger.Error("connect injective db client error:", err)
			return nil, err
		}
		dbClis["injective"] = injectiveDbCli
		valIsJailedChan["injective"] = make(chan []*types.ValIsJail)
		missedSignChan["injective"] = make(chan []*types.ValSignMissed)
		proposalsChan["injective"] = make(chan []*types.Proposal)
		valIsActiveChan["injective"] = make(chan []*types.ValIsActive)
		valRankingChan["injective"] = make(chan []*types.ValRanking)
		preValJailed["injective"] = make(map[string]struct{}, 0)
		preValinActive["injective"] = make(map[string]struct{}, 0)
		preProposalId["injective"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.junoIsMonitored") {
		junoRpcCli, err := junoRpc.InitJunoRpcCli()
		if err != nil {
			logger.Error("connect juno rpc client error: ", err)
			return nil, err
		}
		rpcClis["juno"] = junoRpcCli
		junoDbCli, err := junoDb.InitJunoDbCli()
		if err != nil {
			logger.Error("connect juno db client error:", err)
			return nil, err
		}
		dbClis["juno"] = junoDbCli
		valIsJailedChan["juno"] = make(chan []*types.ValIsJail)
		missedSignChan["juno"] = make(chan []*types.ValSignMissed)
		proposalsChan["juno"] = make(chan []*types.Proposal)
		valIsActiveChan["juno"] = make(chan []*types.ValIsActive)
		valRankingChan["juno"] = make(chan []*types.ValRanking)
		preValJailed["juno"] = make(map[string]struct{}, 0)
		preValinActive["juno"] = make(map[string]struct{}, 0)
		preProposalId["juno"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.neutronconsumerIsMonitored") {
		neutronConsumerRpcCli, err := neutron_consumerRpc.InitNeutronconsumerRpcCli()
		if err != nil {
			logger.Error("connect neutronconsumer rpc client error: ", err)
			return nil, err
		}
		rpcClis["neutronconsumer"] = neutronConsumerRpcCli
		neutronConsumerDbCli, err := neutron_consumerDb.InitNeutronconsumerDbCli()
		if err != nil {
			logger.Error("connect neutronconsumer db client error:", err)
			return nil, err
		}
		dbClis["neutronconsumer"] = neutronConsumerDbCli
		valIsJailedChan["neutronconsumer"] = make(chan []*types.ValIsJail)
		missedSignChan["neutronconsumer"] = make(chan []*types.ValSignMissed)
		proposalsChan["neutronconsumer"] = make(chan []*types.Proposal)
		valIsActiveChan["neutronconsumer"] = make(chan []*types.ValIsActive)
		valRankingChan["neutronconsumer"] = make(chan []*types.ValRanking)
		preValJailed["neutronconsumer"] = make(map[string]struct{}, 0)
		preValinActive["neutronconsumer"] = make(map[string]struct{}, 0)
		preProposalId["neutronconsumer"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.neutronIsMonitored") {
		neutronRpcCli, err := neutronRpc.InitNeutronRpcCli()
		if err != nil {
			logger.Error("connect neutron rpc client error: ", err)
			return nil, err
		}
		rpcClis["neutron"] = neutronRpcCli
		neutronDbCli, err := neutronDb.InitNeutronDbCli()
		if err != nil {
			logger.Error("connect neutron db client error:", err)
			return nil, err
		}
		dbClis["neutron"] = neutronDbCli
		valIsJailedChan["neutron"] = make(chan []*types.ValIsJail)
		missedSignChan["neutron"] = make(chan []*types.ValSignMissed)
		proposalsChan["neutron"] = make(chan []*types.Proposal)
		valIsActiveChan["neutron"] = make(chan []*types.ValIsActive)
		valRankingChan["neutron"] = make(chan []*types.ValRanking)
		preValJailed["neutron"] = make(map[string]struct{}, 0)
		preValinActive["neutron"] = make(map[string]struct{}, 0)
		preProposalId["neutron"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.nyxIsMonitored") {
		nyxRpcCli, err := nyxRpc.InitNyxRpcCli()
		if err != nil {
			logger.Error("connect nyx rpc client error: ", err)
			return nil, err
		}
		rpcClis["nyx"] = nyxRpcCli
		nyxDbCli, err := nyxDb.InitNyxDbCli()
		if err != nil {
			logger.Error("connect nyx db client error:", err)
			return nil, err
		}
		dbClis["nyx"] = nyxDbCli
		valIsJailedChan["nyx"] = make(chan []*types.ValIsJail)
		missedSignChan["nyx"] = make(chan []*types.ValSignMissed)
		proposalsChan["nyx"] = make(chan []*types.Proposal)
		valIsActiveChan["nyx"] = make(chan []*types.ValIsActive)
		valRankingChan["nyx"] = make(chan []*types.ValRanking)
		preValJailed["nyx"] = make(map[string]struct{}, 0)
		preValinActive["nyx"] = make(map[string]struct{}, 0)
		preProposalId["nyx"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.okp4IsMonitored") {
		okp4RpcCli, err := okp4Rpc.InitOkp4RpcCli()
		if err != nil {
			logger.Error("connect okp4 rpc client error: ", err)
			return nil, err
		}
		rpcClis["okp4"] = okp4RpcCli
		okp4DbCli, err := okp4Db.InitOkp4DbCli()
		if err != nil {
			logger.Error("connect okp4 db client error:", err)
			return nil, err
		}
		dbClis["okp4"] = okp4DbCli
		valIsJailedChan["okp4"] = make(chan []*types.ValIsJail)
		missedSignChan["okp4"] = make(chan []*types.ValSignMissed)
		proposalsChan["okp4"] = make(chan []*types.Proposal)
		valIsActiveChan["okp4"] = make(chan []*types.ValIsActive)
		valRankingChan["okp4"] = make(chan []*types.ValRanking)
		preValJailed["okp4"] = make(map[string]struct{}, 0)
		preValinActive["okp4"] = make(map[string]struct{}, 0)
		preProposalId["okp4"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.persistenceIsMonitored") {
		persistenceRpcCli, err := persistenceRpc.InitPersistenceRpcCli()
		if err != nil {
			logger.Error("connect persistence client error: ", err)
			return nil, err
		}
		rpcClis["persistence"] = persistenceRpcCli
		persistenceDbCli, err := persistenceDb.InitPersistenceDbCli()
		if err != nil {
			logger.Error("connect persistence db client error:", err)
			return nil, err
		}
		dbClis["persistence"] = persistenceDbCli
		valIsJailedChan["persistence"] = make(chan []*types.ValIsJail)
		missedSignChan["persistence"] = make(chan []*types.ValSignMissed)
		proposalsChan["persistence"] = make(chan []*types.Proposal)
		valIsActiveChan["persistence"] = make(chan []*types.ValIsActive)
		valRankingChan["persistence"] = make(chan []*types.ValRanking)
		preValJailed["persistence"] = make(map[string]struct{}, 0)
		preValinActive["persistence"] = make(map[string]struct{}, 0)
		preProposalId["persistence"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.providerIsMonitored") {
		providerRpcCli, err := providerRpc.InitProviderRpcCli()
		if err != nil {
			logger.Error("connect provider rpc client error: ", err)
			return nil, err
		}
		rpcClis["provider"] = providerRpcCli
		providerDbCli, err := providerDb.InitProviderDbCli()
		if err != nil {
			logger.Error("connect provider db client error:", err)
			return nil, err
		}
		dbClis["provider"] = providerDbCli
		valIsJailedChan["provider"] = make(chan []*types.ValIsJail)
		missedSignChan["provider"] = make(chan []*types.ValSignMissed)
		proposalsChan["provider"] = make(chan []*types.Proposal)
		valIsActiveChan["provider"] = make(chan []*types.ValIsActive)
		valRankingChan["provider"] = make(chan []*types.ValRanking)
		preValJailed["provider"] = make(map[string]struct{}, 0)
		preValinActive["provider"] = make(map[string]struct{}, 0)
		preProposalId["provider"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.rizonIsMonitored") {
		rizonRpcCli, err := rizonRpc.InitCosmosRpcCli()
		if err != nil {
			logger.Error("connect rizon rpc client error: ", err)
			return nil, err
		}
		rpcClis["rizon"] = rizonRpcCli
		rizonDbCli, err := rizonDb.InitRizonDbCli()
		if err != nil {
			logger.Error("connect rizon db client error:", err)
			return nil, err
		}
		dbClis["rizon"] = rizonDbCli
		valIsJailedChan["rizon"] = make(chan []*types.ValIsJail)
		missedSignChan["rizon"] = make(chan []*types.ValSignMissed)
		proposalsChan["rizon"] = make(chan []*types.Proposal)
		valIsActiveChan["rizon"] = make(chan []*types.ValIsActive)
		valRankingChan["rizon"] = make(chan []*types.ValRanking)
		preValJailed["rizon"] = make(map[string]struct{}, 0)
		preValinActive["rizon"] = make(map[string]struct{}, 0)
		preProposalId["rizon"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.secretIsMonitored") {
		secretRpcCli, err := secretRpc.InitCosmosRpcCli()
		if err != nil {
			logger.Error("connect secret rpc client error: ", err)
			return nil, err
		}
		rpcClis["secret"] = secretRpcCli
		secretDbCli, err := secretDb.InitSecretDbCli()
		if err != nil {
			logger.Error("connect secret db client error:", err)
			return nil, err
		}
		dbClis["secret"] = secretDbCli
		valIsJailedChan["secret"] = make(chan []*types.ValIsJail)
		missedSignChan["secret"] = make(chan []*types.ValSignMissed)
		proposalsChan["secret"] = make(chan []*types.Proposal)
		valIsActiveChan["secret"] = make(chan []*types.ValIsActive)
		valRankingChan["secret"] = make(chan []*types.ValRanking)
		preValJailed["secret"] = make(map[string]struct{}, 0)
		preValinActive["secret"] = make(map[string]struct{}, 0)
		preProposalId["secret"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.sommelierIsMonitored") {
		sommelierRpcCli, err := sommelierRpc.InitSommelierRpcCli()
		if err != nil {
			logger.Error("connect sommelier rpc client error: ", err)
			return nil, err
		}
		rpcClis["sommelier"] = sommelierRpcCli
		sommelierDbCli, err := sommelierDb.InitSommelierDbCli()
		if err != nil {
			logger.Error("connect sommelier db client error:", err)
			return nil, err
		}
		dbClis["sommelier"] = sommelierDbCli
		valIsJailedChan["sommelier"] = make(chan []*types.ValIsJail)
		missedSignChan["sommelier"] = make(chan []*types.ValSignMissed)
		proposalsChan["sommelier"] = make(chan []*types.Proposal)
		valIsActiveChan["sommelier"] = make(chan []*types.ValIsActive)
		valRankingChan["sommelier"] = make(chan []*types.ValRanking)
		preValJailed["sommelier"] = make(map[string]struct{}, 0)
		preValinActive["sommelier"] = make(map[string]struct{}, 0)
		preProposalId["sommelier"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.sputnikIsMonitored") {
		sputnikRpcCli, err := sputnikRpc.InitSputnikRpcCli()
		if err != nil {
			logger.Error("connect sputnik rpc client error: ", err)
			return nil, err
		}
		rpcClis["sputnik"] = sputnikRpcCli
		sputnikDbCli, err := sputnikDb.InitSputnikDbCli()
		if err != nil {
			logger.Error("connect sputnik db client error:", err)
			return nil, err
		}
		dbClis["sputnik"] = sputnikDbCli
		valIsJailedChan["sputnik"] = make(chan []*types.ValIsJail)
		missedSignChan["sputnik"] = make(chan []*types.ValSignMissed)
		proposalsChan["sputnik"] = make(chan []*types.Proposal)
		valIsActiveChan["sputnik"] = make(chan []*types.ValIsActive)
		valRankingChan["sputnik"] = make(chan []*types.ValRanking)
		preValJailed["sputnik"] = make(map[string]struct{}, 0)
		preValinActive["sputnik"] = make(map[string]struct{}, 0)
		preProposalId["sputnik"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.teritoriIsMonitored") {
		teritoriRpcCli, err := teritoriRpc.InitTeritoriRpcCli()
		if err != nil {
			logger.Error("connect teritori rpc client error: ", err)
			return nil, err
		}
		rpcClis["teritori"] = teritoriRpcCli
		teritoriDbCli, err := teritoriDb.InitTeritoriDbCli()
		if err != nil {
			logger.Error("connect teritori db client error:", err)
			return nil, err
		}
		dbClis["teritori"] = teritoriDbCli
		valIsJailedChan["teritori"] = make(chan []*types.ValIsJail)
		missedSignChan["teritori"] = make(chan []*types.ValSignMissed)
		proposalsChan["teritori"] = make(chan []*types.Proposal)
		valIsActiveChan["teritori"] = make(chan []*types.ValIsActive)
		valRankingChan["teritori"] = make(chan []*types.ValRanking)
		preValJailed["teritori"] = make(map[string]struct{}, 0)
		preValinActive["teritori"] = make(map[string]struct{}, 0)
		preProposalId["teritori"] = make(map[int64]struct{}, 0)
	}

	if viper.GetBool("alert.xplaIsMonitored") {
		xplaRpcCli, err := xplaRpc.InitXplaRpcCli()
		if err != nil {
			logger.Error("connect xpla rpc client error: ", err)
			return nil, err
		}
		rpcClis["xpla"] = xplaRpcCli
		xplaDbCli, err := xplaDb.InitXplaDbCli()
		if err != nil {
			logger.Error("connect xpla db client error:", err)
			return nil, err
		}
		dbClis["xpla"] = xplaDbCli
		valIsJailedChan["xpla"] = make(chan []*types.ValIsJail)
		missedSignChan["xpla"] = make(chan []*types.ValSignMissed)
		proposalsChan["xpla"] = make(chan []*types.Proposal)
		valIsActiveChan["xpla"] = make(chan []*types.ValIsActive)
		valRankingChan["xpla"] = make(chan []*types.ValRanking)
		preValJailed["xpla"] = make(map[string]struct{}, 0)
		preValinActive["xpla"] = make(map[string]struct{}, 0)
		preProposalId["xpla"] = make(map[int64]struct{}, 0)
	}

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
		valIsJailedChan: valIsJailedChan,
		missedSignChan:  missedSignChan,
		proposalsChan:   proposalsChan,
		valIsActiveChan: valIsActiveChan,
		valRankingChan:  valRankingChan,
	}, nil
}

func (m *Monitor) Start() {
	epochTicker := time.NewTicker(time.Duration(viper.GetInt("alert.timeInterval")) * time.Second)
	consumerChains := make(map[string]string, 0)
	consumerChains["apollo"] = "provider"
	consumerChains["gopher"] = "provider"
	consumerChains["hero"] = "provider"
	consumerChains["neutronconsumer"] = "provider"
	consumerChains["sputnik"] = "provider"

	for project := range m.RpcClis {
		startHeight[project] = m.RpcClis[project].GetBlockHeight()
	}
	for range epochTicker.C {
		for project := range m.RpcClis {
			if project == "apollo" || project == "gopher" || project == "hero" || project == "neutronconsumer" || project == "sputnik" {
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

	mos := make(map[string]*types.MonitorObj)
	for _, moDb := range mosDb {
		mos[moDb.OperatorAddr] = moDb
	}
	operatorAddrs := config.GetOperatorAddrs(consumerChain)
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
	operatorAddrs := config.GetOperatorAddrs(project)

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
		newpreValJailed := make(map[string]map[string]struct{}, 0)
		newpreValJailed[caredData.ChainName] = make(map[string]struct{}, 0)
		newpreValinActive := make(map[string]map[string]struct{}, 0)
		newpreValinActive[caredData.ChainName] = make(map[string]struct{}, 0)

		for _, valInfo := range caredData.ValInfos {
			if _, ok := preValJailed[caredData.ChainName][valInfo.Moniker]; !ok && valInfo.Jailed {
				valIsJailed = append(valIsJailed, &types.ValIsJail{
					ChainName: caredData.ChainName,
					Moniker:   valInfo.Moniker,
				})
			}
			if _, ok := preValinActive[caredData.ChainName][valInfo.Moniker]; !ok && valInfo.Status != 3 {
				valIsActive = append(valIsActive, &types.ValIsActive{
					ChainName: caredData.ChainName,
					Moniker:   valInfo.Moniker,
					Status:    valInfo.Status,
				})
			}
			if valInfo.Jailed {
				newpreValJailed[caredData.ChainName][valInfo.Moniker] = struct{}{}
			}
			if valInfo.Status != 3 {
				newpreValinActive[caredData.ChainName][valInfo.Moniker] = struct{}{}
			}
		}
		preValJailed[caredData.ChainName] = newpreValJailed[caredData.ChainName]
		preValinActive[caredData.ChainName] = newpreValinActive[caredData.ChainName]

		m.valIsJailedChan[caredData.ChainName] <- valIsJailed
		m.valIsActiveChan[caredData.ChainName] <- valIsActive
	}

	if caredData.Proposals != nil && len(caredData.Proposals) > 0 {
		logger.Info("Start saving proposals information")
		err := m.DbClis[caredData.ChainName].BatchSaveProposals(caredData.Proposals)
		if err != nil {
			logger.Error("save the proposals information fail:", err)
		}
		logger.Info("Save the proposals information successfully")

		proposals := make([]*types.Proposal, 0)
		newPreProposalId := make(map[string]map[int64]struct{}, 0)
		newPreProposalId[caredData.ChainName] = make(map[int64]struct{}, 0)
		for _, proposal := range caredData.Proposals {
			if _, ok := newPreProposalId[caredData.ChainName][proposal.ProposalId]; !ok {
				newPreProposalId[caredData.ChainName][proposal.ProposalId] = struct{}{}
			}

			if _, ok := preProposalId[caredData.ChainName][proposal.ProposalId]; !ok {
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
				preProposalId[caredData.ChainName][proposal.ProposalId] = struct{}{}
			}
		}
		preProposalId[caredData.ChainName] = newPreProposalId[caredData.ChainName]
		m.proposalsChan[caredData.ChainName] <- proposals
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
		m.missedSignChan[caredData.ChainName] <- missedSign
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
		m.valRankingChan[caredData.ChainName] <- valRankings
	}

	timeInterval := viper.GetInt("alert.timeInterval")
	endHeight := end / int64(timeInterval) * int64(timeInterval)
	if monitorHeight[caredData.ChainName] != endHeight {
		m.DbClis[caredData.ChainName].BatchSaveValStats(endHeight-int64(timeInterval)+int64(1), endHeight)
		monitorHeight[caredData.ChainName] = endHeight
	}
}

func (m *Monitor) SendEmail() {
	time.Sleep(time.Second * 10)
	for project := range m.RpcClis {
		mailSender := viper.GetString("mail.sender")
		receiver1Conf := fmt.Sprintf("mail.%sReceiver1", project)
		receiver2Conf := fmt.Sprintf("mail.%sReceiver2", project)
		receiver1 := viper.GetString(receiver1Conf)
		receiver2 := viper.GetString(receiver2Conf)
		mailReceiver := strings.Join([]string{receiver1, receiver2}, ",")
		go m.sendEmail(mailSender, receiver1, mailReceiver, project)
	}
}

func (m *Monitor) sendEmail(mailSender, receiver1, mailReceiver, project string) {
	for {
		select {
		case valJailed := <-m.valIsJailedChan[project]:
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
		case valisAtive := <-m.valIsActiveChan[project]:
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
		case missedSign := <-m.missedSignChan[project]:
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
		case proposals := <-m.proposalsChan[project]:
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
		case valRanking := <-m.valRankingChan[project]:
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
