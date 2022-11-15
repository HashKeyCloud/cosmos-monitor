package cosmos_rpc

import (
	"fmt"
	"testing"
	"time"

	"cosmosmonitor/rpc"
	"cosmosmonitor/types"
	base "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	distribution "cosmossdk.io/api/cosmos/distribution/v1beta1"
	gov "cosmossdk.io/api/cosmos/gov/v1beta1"
	staking "cosmossdk.io/api/cosmos/staking/v1beta1"
)

func TestGovInfo(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("cosmos-grpc.polkachu.com:14990")
	if err != nil {
		logger.Error("Failed to create cosmos gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &CosmosCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	monitorObj := make([]*types.MonitorObj, 0)
	monitorObj = append(monitorObj, &types.MonitorObj{
		"Coinbase Custody",
		"cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en",
		"d68eec0d2e8248f1ec64cdb585edb61eca432bd8",
		"cosmos1c4k24jzduc365kywrsvf5ujz4ya6mwymy8vq4q",
	})
	proposals, _ := cc.GetProposal(monitorObj)
	for _, p := range proposals {
		fmt.Println(p)
	}

}

func TestGetValInfo(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("cosmos-grpc.polkachu.com:14990")
	if err != nil {
		logger.Error("Failed to create cosmos gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &CosmosCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "cosmosvaloper12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5gk7zg2")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("cosmos-grpc.polkachu.com:14990")
	if err != nil {
		logger.Error("Failed to create cosmos gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &CosmosCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObj := make([]*types.MonitorObj, 0)
	m1 := &types.MonitorObj{
		"Coinbase Custody",
		"cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en",
		"d68eec0d2e8248f1ec64cdb585edb61eca432bd8",
		"cosmos1c4k24jzduc365kywrsvf5ujz4ya6mwymy8vq4q",
	}
	monitorObj = append(monitorObj, m1)
	proposalAssignments, signs, signsmissed, _ := cc.GetValPerformance(12775580, monitorObj)
	for _, proposalAssignment := range proposalAssignments {
		fmt.Println("proposalAssignment:", proposalAssignment)
	}
	for _, sign := range signs {
		fmt.Println("sign:", sign)
	}

	for _, signmissed := range signsmissed {
		fmt.Println("sign missed:", signmissed)
	}
}

func TestGetValRanking(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("cosmos-grpc.polkachu.com:14990")
	if err != nil {
		logger.Error("Failed to create cosmos gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &CosmosCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObj := make([]*types.MonitorObj, 0)
	m1 := &types.MonitorObj{
		"Coinbase Custody",
		"cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en",
		"d68eec0d2e8248f1ec64cdb585edb61eca432bd8",
		"cosmos1c4k24jzduc365kywrsvf5ujz4ya6mwymy8vq4q",
	}
	/*m1 := &types.MonitorObj{
		"Huobi-1",
		"cosmosvaloper12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5gk7zg2",
		"4af69d6a5436c30e3584c1628433de55e758bcca",
		"cosmos12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5dz2hye",
	}*/
	monitorObj = append(monitorObj, m1)
	valRanking, err := cc.GetValRanking(monitorObj, "cosmos")
	if err != nil {
		logger.Error("Failed to query LatestValidatorSet, err:", err)
	}
	for _, v := range valRanking {
		logger.Info("validator ranking:", v)
	}
}
