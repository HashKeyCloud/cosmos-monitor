package rizon_rpc

import (
	"fmt"
	"testing"

	"cosmosmonitor/rpc"
	"cosmosmonitor/types"
	base "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	distribution "cosmossdk.io/api/cosmos/distribution/v1beta1"
	gov "cosmossdk.io/api/cosmos/gov/v1beta1"
	staking "cosmossdk.io/api/cosmos/staking/v1beta1"
)

func TestGetValInfo(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("https://rizon.nodejumper.io:9090")
	if err != nil {
		logger.Error("Failed to create rizon gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &RizonCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}
	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "rizonvaloper1qpgalplval6ynvq2r0afxchmc0x00zmuu06300")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetProposal(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("persistence-grpc.polkachu.com:15490")
	if err != nil {
		logger.Error("Failed to create rizon gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &RizonCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObjs := make([]*types.MonitorObj, 0)
	monitorObjs = append(monitorObjs, &types.MonitorObj{
		Moniker:         "Figment",
		OperatorAddr:    "nvaloper1e22wh85arrkpe6derct4l5r4gp3awvase24pf6",
		OperatorAddrHex: "e2e65903bab5344e5cffddc2cd638a6bab2e2e79",
		SelfStakeAddr:   "n1e22wh85arrkpe6derct4l5r4gp3awvas0rajn8",
	})
	monitors, _ := cc.GetProposal(monitorObjs)
	for _, monitor := range monitors {
		fmt.Println("proposal:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("persistence-grpc.polkachu.com:15490")
	if err != nil {
		logger.Error("Failed to create rizon gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &RizonCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObjs := make([]*types.MonitorObj, 0)
	monitorObjs = append(monitorObjs, &types.MonitorObj{
		Moniker:         "Cosmostation",
		OperatorAddr:    "junovaloper1t8ehvswxjfn3ejzkjtntcyrqwvmvuknzmvtaaa",
		OperatorAddrHex: "80f24bfda3e6a8c1bac0517e7665ac9145d609f7",
		SelfStakeAddr:   "juno1t8ehvswxjfn3ejzkjtntcyrqwvmvuknzy3ajxy",
	})
	proposalAssignment, valSign, valSignMissed, _ := cc.GetValPerformance(4299893, monitorObjs)
	for _, monitor := range proposalAssignment {
		fmt.Println("proposalAssignment:", monitor)
	}

	for _, monitor := range valSign {
		fmt.Println("valSign:", monitor)
	}

	for _, monitor := range valSignMissed {
		fmt.Println("valSignMissed:", monitor)
	}
}
