package cosmos_rpc

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

type CosmosCli struct {
	*rpc.ChainCli
}

// NewRpcCli Create a new RPC service
func NewCosmosRpcCli() (*CosmosCli, error) {
	endpoint := fmt.Sprintf("%s:%s", viper.GetString("cosmos.ip"), viper.GetString("cosmos.gRPCport"))
	grpcConn, err := rpc.NewChainRpcCli(endpoint)
	if err != nil {
		logger.Error("Failed to create cosmos gRPC client, err:", err)
		return nil, err
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	cli := CosmosCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}
	return &cli, err
}

var logger = log.RPCLogger.WithField("module", "rpc")