package neutron_rpc

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

type NeutronCli struct {
	*rpc.ChainCli
}

func InitNeutronRpcCli() (*NeutronCli, error) {
	endpoint := fmt.Sprintf("%s:%s", viper.GetString("gRpc.neutronIp"), viper.GetString("gRpc.neutronPort"))
	grpcConn, err := rpc.InitChainRpcCli(endpoint)
	if err != nil {
		logger.Error("Failed to create neutron gRPC client, err:", err)
		return nil, err
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)
	govQueryCli := gov.NewQueryClient(grpcConn)
	baseCli := base.NewServiceClient(grpcConn)

	return &NeutronCli{
		ChainCli: &rpc.ChainCli{
			StakingQueryCli: stakingQueryCli,
			DistributionCli: distributionCli,
			GovQueryCli:     govQueryCli,
			BaseQuaryCli:    baseCli,
		},
	}, nil
}

var logger = log.RPCLogger.WithField("module", "rpc")
