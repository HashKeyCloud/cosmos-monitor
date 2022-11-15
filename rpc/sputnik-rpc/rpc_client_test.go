package sputnik_rpc

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

func TestGetValPerformance(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("xxx")
	if err != nil {
		logger.Error("Failed to create apollo gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &SputnikCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObj := make([]*types.MonitorObj, 0)
	m1 := &types.MonitorObj{
		"HashQuark",
		"cosmosvaloper10e4rluum506yc63vgrdcn500zyev2c5jw2ndkv",
		"7B0DE97C082AE7911928004FE8B145AA1FF56038",
		"cosmos10e4rluum506yc63vgrdcn500zyev2c5jt78c6l",
	}
	monitorObj = append(monitorObj, m1)
	proposalAssignments, signs, signsmissed, _ := cc.GetValPerformance(46950, monitorObj)
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
