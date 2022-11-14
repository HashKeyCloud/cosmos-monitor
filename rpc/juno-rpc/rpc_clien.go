package juno_rpc

import (
	"cosmosmonitor/log"
	"cosmosmonitor/rpc"
	base "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	distribution "cosmossdk.io/api/cosmos/distribution/v1beta1"
	gov "cosmossdk.io/api/cosmos/gov/v1beta1"
	staking "cosmossdk.io/api/cosmos/staking/v1beta1"
	"fmt"
	"github.com/spf13/viper"
)

type JunoCli struct {
	*rpc.ChainCli
}

func InitJunoRpcCli() (*JunoCli, error) {
	endpoint := fmt.Sprintf("%s:%s", viper.GetString("gRpc.junoIp"), viper.GetString("gRpc.junoPort"))
	grpcConn, err := rpc.InitChainRpcCli(endpoint)
	if err != nil {
		logger.Error("Failed to create juno gRPC client, err:", err)
		return nil, err
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)

	return &JunoCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			DistributionCli: distributionCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
		},
	}, nil
}

var logger = log.RPCLogger.WithField("module", "rpc")
