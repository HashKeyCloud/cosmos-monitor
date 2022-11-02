package injective_rpc

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

type InjectiveCli struct {
	*rpc.ChainCli
}

func NewInjectiveRpcCli() (*InjectiveCli, error) {
	endpoint := fmt.Sprintf("%s:%s", viper.GetString("injective.ip"), viper.GetString("injective.gRPCport"))
	grpcConn, err := rpc.NewChainRpcCli(endpoint)
	if err != nil {
		logger.Error("Failed to create injective gRPC client, err:", err)
		return nil, err
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)

	return &InjectiveCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			DistributionCli: distributionCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
		},
	}, nil
}

var logger = log.RPCLogger.WithField("module", "rpc")
