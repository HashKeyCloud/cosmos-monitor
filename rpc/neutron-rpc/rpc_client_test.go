package neutron_rpc

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
	grpcConn, err := rpc.InitChainRpcCli("xxxx")
	if err != nil {
		logger.Error("Failed to create neutron gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &NeutronCli{
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
		"Neutron Validator",
		"neutronvaloper1w5a2fmp3e6fjxes8rc36fneflvmq3asez7g6kl",
		"aa5c70f5a3f259f4186f8a8e704cd992f6635b3a",
		"neutron1w5a2fmp3e6fjxes8rc36fneflvmq3asecr3rsm",
	})
	proposals, _ := cc.GetProposal(monitorObj)
	for _, p := range proposals {
		fmt.Println(p)
	}

}

func TestGetValInfo(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("xxx")
	if err != nil {
		logger.Error("Failed to create neutron gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &NeutronCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "neutronvaloper1w5a2fmp3e6fjxes8rc36fneflvmq3asez7g6kl")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	grpcConn, err := rpc.InitChainRpcCli("xxx")
	if err != nil {
		logger.Error("Failed to create neutron gRPC client, err:", err)
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cc := &NeutronCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}

	monitorObj := make([]*types.MonitorObj, 0)
	m1 := &types.MonitorObj{
		"Neutron Validator",
		"neutronvaloper1w5a2fmp3e6fjxes8rc36fneflvmq3asez7g6kl",
		"aa5c70f5a3f259f4186f8a8e704cd992f6635b3a",
		"neutron1w5a2fmp3e6fjxes8rc36fneflvmq3asecr3rsm",
	}
	monitorObj = append(monitorObj, m1)
	proposalAssignments, signs, signsmissed, _ := cc.GetValPerformance(102644, monitorObj)
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
