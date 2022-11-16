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
	grpcConn, err := rpc.InitChainRpcCli("xx")
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
	monitorObj = append(monitorObj, "neutronvaloper1d88mk88kvged9z58xw5h89wke65lu8vmd02rga")
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

func TestGetValRanking(t *testing.T) {
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

	monitorObj := make([]*types.MonitorObj, 0)
	m1 := &types.MonitorObj{
		"HashQuark",
		"neutronvaloper1d88mk88kvged9z58xw5h89wke65lu8vmd02rga",
		"b6dd95a54dec1130e86a57798d87d04fbbc9f4d9",
		"neutron1d88mk88kvged9z58xw5h89wke65lu8vmhjn6we",
	}
	/*m1 := &types.MonitorObj{
		"Huobi-1",
		"cosmosvaloper12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5gk7zg2",
		"4af69d6a5436c30e3584c1628433de55e758bcca",
		"cosmos12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5dz2hye",
	}*/
	monitorObj = append(monitorObj, m1)
	valRanking, err := cc.GetValRanking(monitorObj, "neutron")
	if err != nil {
		logger.Error("Failed to query LatestValidatorSet, err:", err)
	}
	for _, v := range valRanking {
		logger.Info("validator ranking:", v)
	}
}
