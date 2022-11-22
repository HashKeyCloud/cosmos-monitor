package akash_rpc

import (
	"fmt"

	"github.com/spf13/viper"

	"cosmosmonitor/log"
	"cosmosmonitor/rpc"
	base "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	distribution "cosmossdk.io/api/cosmos/distribution/v1beta1"
	gov "cosmossdk.io/api/cosmos/gov/v1beta1"
	staking "cosmossdk.io/api/cosmos/staking/v1beta1"
)

type AkashCli struct {
	*rpc.ChainCli
}

// InitRpcCli Create a new RPC service
func InitAkashRpcCli() (*AkashCli, error) {
	endpoint := fmt.Sprintf("%s:%s", viper.GetString("gRpc.akashIp"), viper.GetString("gRpc.akashPort"))
	grpcConn, err := rpc.InitChainRpcCli(endpoint)
	if err != nil {
		logger.Error("Failed to create akash gRPC client, err:", err)
		return nil, err
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	return &AkashCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}, err
}

var logger = log.RPCLogger.WithField("module", "rpc")