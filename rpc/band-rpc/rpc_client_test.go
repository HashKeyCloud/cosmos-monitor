package band_rpc

import (
	"cosmosmonitor/rpc"
	"cosmosmonitor/types"
	base "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	distribution "cosmossdk.io/api/cosmos/distribution/v1beta1"
	gov "cosmossdk.io/api/cosmos/gov/v1beta1"
	staking "cosmossdk.io/api/cosmos/staking/v1beta1"
	"fmt"
	"testing"
	"time"
)

func TestGovInfo(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("xxx")
	if err != nil {
		logger.Error("Failed to create band gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &BandCli{
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
	grpcConn, err := rpc.InitChainRpcCli("xxxxx")
	if err != nil {
		logger.Error("Failed to create band gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &BandCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "xxxx")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("xxxxx")
	if err != nil {
		logger.Error("Failed to create band gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &BandCli{
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
