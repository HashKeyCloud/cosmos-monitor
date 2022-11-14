package teritori_rpc

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
	grpcConn, err := rpc.InitChainRpcCli("https://grpc.teritori.silknodes.io/")
	if err != nil {
		logger.Error("Failed to create teritori gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &TeritoriCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}
	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "torivaloper1xu736l4vt6l2pg9k2yk66fq7zq6y4aj5xmd6vq")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetProposal(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("teritori-grpc.polkachu.com:15990")
	if err != nil {
		logger.Error("Failed to create sommelier gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &TeritoriCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObjs := make([]*types.MonitorObj, 0)
	monitorObjs = append(monitorObjs, &types.MonitorObj{
		Moniker:         "Polychain",
		OperatorAddr:    "sommvaloper1a583xjv6ylrdfzm3udk6jnj4xy28k2wge0eazz",
		OperatorAddrHex: "f49ff6e3b6a00602edb46e1a1f8fbe79f4880ce7",
		SelfStakeAddr:   "somm1a583xjv6ylrdfzm3udk6jnj4xy28k2wgv3hq9d",
	})
	monitors, _ := cc.GetProposal(monitorObjs)
	for _, monitor := range monitors {
		fmt.Println("proposal:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("teritori-grpc.polkachu.com:15990")
	if err != nil {
		logger.Error("Failed to create sommelier gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &TeritoriCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObjs := make([]*types.MonitorObj, 0)
	monitorObjs = append(monitorObjs, &types.MonitorObj{
		Moniker:         "Polychain",
		OperatorAddr:    "sommvaloper1a583xjv6ylrdfzm3udk6jnj4xy28k2wge0eazz",
		OperatorAddrHex: "f49ff6e3b6a00602edb46e1a1f8fbe79f4880ce7",
		SelfStakeAddr:   "somm1a583xjv6ylrdfzm3udk6jnj4xy28k2wgv3hq9d",
	})
	proposalAssignment, valSign, valSignMissed, _ := cc.GetValPerformance(6275552, monitorObjs)
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
