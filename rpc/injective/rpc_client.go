package injective

import (
	"context"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"cosmosmonitor/utils"
	distribution "cosmossdk.io/api/cosmos/distribution/v1beta1"
	staking "cosmossdk.io/api/cosmos/staking/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InjectiveCli struct {
	GRPCCli         grpc.ClientConn
	StakingQueryCli staking.QueryClient
	DistributionCli distribution.QueryClient
}

func NewInjectiveRpcCli(endpoint string) (*InjectiveCli, error) {
	dialOpts := []grpc.DialOption{
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Maximum receive value 128 MB
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(128 * 1024 * 1024)),
	}
	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(endpoint, dialOpts...)
	if err != nil {
		logger.Error("Failed to create cosmos gRPC client, err:", err)
		return nil, err
	}

	stakingQueryCli := staking.NewQueryClient(grpcConn)
	distributionCli := distribution.NewQueryClient(grpcConn)
	return &InjectiveCli{
		GRPCCli:         *grpcConn,
		StakingQueryCli: stakingQueryCli,
		DistributionCli: distributionCli,
	}, nil
}

// GetValInfo Get validator information via RPC
func (ic *InjectiveCli) GetValInfo(operatorAddrs []string) ([]*types.ValInfo, error) {
	valsInfo := make([]*types.ValInfo, 0)
	for _, operatorAddr := range operatorAddrs {
		qVal := &staking.QueryValidatorRequest{
			ValidatorAddr: operatorAddr,
		}
		queryValidatorResponse, err := ic.StakingQueryCli.Validator(context.Background(), qVal)
		if err != nil {
			logger.Error("Failed to query validator information, err:", err)
			return nil, err
		}
		selfAddr := utils.Operator2SelfAddrInj(operatorAddr)
		operatorAddrHex := utils.Puk2addrHex(queryValidatorResponse.Validator.ConsensusPubkey.GetValue())
		dVal := &distribution.QueryDelegatorWithdrawAddressRequest{
			DelegatorAddress: selfAddr,
		}
		rewardAddr, err := ic.DistributionCli.DelegatorWithdrawAddress(context.Background(), dVal)
		if err != nil {
			logger.Error("Failed to query ValidatorDistributionInfo, err:", err)
			return nil, err
		}

		commissionRatesfloat := utils.Div18(queryValidatorResponse.Validator.Commission.CommissionRates.Rate)
		maxRatefloat := utils.Div18(queryValidatorResponse.Validator.Commission.CommissionRates.MaxRate)
		maxChangeRatefloat := utils.Div18(queryValidatorResponse.Validator.Commission.CommissionRates.MaxChangeRate)

		valInfo := &types.ValInfo{
			Moniker:           queryValidatorResponse.Validator.Description.Moniker,
			OperatorAddr:      operatorAddr,
			OperatorAddrHex:   operatorAddrHex,
			SelfStakeAddr:     selfAddr,
			RewardAddr:        rewardAddr.GetWithdrawAddress(),
			Jailed:            queryValidatorResponse.Validator.GetJailed(),
			Status:            int32(queryValidatorResponse.Validator.GetStatus()),
			VotingPower:       queryValidatorResponse.Validator.Tokens,
			Identity:          queryValidatorResponse.Validator.Description.Identity,
			Website:           queryValidatorResponse.Validator.Description.Website,
			Details:           queryValidatorResponse.Validator.Description.Details,
			SecurityContact:   queryValidatorResponse.Validator.Description.SecurityContact,
			CommissionRates:   commissionRatesfloat,
			MaxRate:           maxRatefloat,
			MaxChangeRate:     maxChangeRatefloat,
			MinSelfDelegation: queryValidatorResponse.Validator.MinSelfDelegation,
		}
		valsInfo = append(valsInfo, valInfo)
	}

	return valsInfo, nil
}

var logger = log.RPCLogger.WithField("module", "rpc")
