package zeta_rpc

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

type ZetaCli struct {
	*rpc.ChainCli
}

// InitRpcCli Create a new RPC service
func InitZetaRpcCli() (*ZetaCli, error) {
	endpoint := fmt.Sprintf("%s:%s", viper.GetString("gRpc.zetaIp"), viper.GetString("gRpc.zetaPort"))
	grpcConn, err := rpc.InitChainRpcCli(endpoint)
	if err != nil {
		logger.Error("Failed to create zeta gRPC client, err:", err)
		return nil, err
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)

	return &ZetaCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
			DistributionCli: distributionCli,
		},
	}, err
}

var logger = log.RPCLogger.WithField("module", "rpc")
