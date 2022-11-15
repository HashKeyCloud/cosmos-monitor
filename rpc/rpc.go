package rpc

import (
	"context"
	"encoding/hex"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"cosmosmonitor/utils"
	base "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	distribution "cosmossdk.io/api/cosmos/distribution/v1beta1"
	gov "cosmossdk.io/api/cosmos/gov/v1beta1"
	staking "cosmossdk.io/api/cosmos/staking/v1beta1"
	blockTypes "cosmossdk.io/api/tendermint/types"
)

type Client interface {
	GetValInfo(operatorAddrs []string) ([]*types.ValInfo, error)
	GetProposal(monitorObjs []*types.MonitorObj) ([]*types.Proposal, error)
	GetValPerformance(start int64, monitorObjs []*types.MonitorObj) ([]*types.ProposalAssignment, []*types.ValSign, []*types.ValSignMissed, error)
	GetValRanking(monitorObjs []*types.MonitorObj, project string) ([]*types.ValRanking, error)
}

type ChainCli struct {
	StakingQueryCli staking.QueryClient
	GovQueryCli     gov.QueryClient
	BaseQuaryCli    base.ServiceClient
	DistributionCli distribution.QueryClient
}

func InitChainRpcCli(endpoint string) (*grpc.ClientConn, error) {
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
	return grpcConn, nil
}

// GetValInfo Get validator information via RPC
func (cc *ChainCli) GetValInfo(operatorAddrs []string) ([]*types.ValInfo, error) {
	valsInfo := make([]*types.ValInfo, 0)
	for _, operatorAddr := range operatorAddrs {
		qVal := &staking.QueryValidatorRequest{
			ValidatorAddr: operatorAddr,
		}
		queryValidatorResponse, err := cc.StakingQueryCli.Validator(context.Background(), qVal)
		if err != nil {
			logger.Error("Failed to query validator information, err:", err)
			return nil, err
		}
		selfAddr := utils.Operator2SelfAddr(operatorAddr)
		operatorAddrHex := utils.Puk2addrHex(queryValidatorResponse.Validator.ConsensusPubkey.GetValue())
		dVal := &distribution.QueryDelegatorWithdrawAddressRequest{
			DelegatorAddress: selfAddr,
		}
		rewardAddr, err := cc.DistributionCli.DelegatorWithdrawAddress(context.Background(), dVal)
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

// GetProposal Obtain proposals in VOTING PERIOD through RPC, and obtain validator votes.
func (cc *ChainCli) GetProposal(monitorObjs []*types.MonitorObj) ([]*types.Proposal, error) {
	queryProposalsRequest := &gov.QueryProposalsRequest{}
	queryProposalsRespons, err := cc.GovQueryCli.Proposals(context.Background(), queryProposalsRequest)
	if err != nil {
		logger.Error("Failed to query proposal, err:", err)
		return nil, err
	}
	proposals := make([]*types.Proposal, 0)
	for _, proposal := range queryProposalsRespons.Proposals {
		if proposal.Status == gov.ProposalStatus(2) {
			for _, monitorObj := range monitorObjs {
				queryVoteRequest := &gov.QueryVoteRequest{
					ProposalId: proposal.ProposalId,
					Voter:      monitorObj.SelfStakeAddr,
				}

				queryVoteResponse, err := cc.GovQueryCli.Vote(context.Background(), queryVoteRequest)

				var status int32
				if err != nil && strings.Contains(err.Error(), "code = InvalidArgument desc") {
					status = 0
				} else if err == nil {
					status = int32(queryVoteResponse.Vote.Options[len(queryVoteResponse.Vote.Options)-1].Option)
				} else {
					return nil, err
				}

				proposals = append(proposals, &types.Proposal{
					ProposalId:      int64(proposal.ProposalId),
					VotingStartTime: proposal.VotingStartTime.AsTime().Format("2006-01-02 15:04:05"),
					VotingEndTime:   proposal.VotingEndTime.AsTime().Format("2006-01-02 15:04:05"),
					Description:     proposal.GetContent().String(),
					Moniker:         monitorObj.Moniker,
					OperatorAddr:    monitorObj.OperatorAddr,
					Status:          status,
				})
			}
		}
	}
	return proposals, nil
}

// GetValPerformance Query the validator's block production and signature status through RPC
func (cc *ChainCli) GetValPerformance(start int64, monitorObjs []*types.MonitorObj) ([]*types.ProposalAssignment, []*types.ValSign, []*types.ValSignMissed, error) {
	blockByHeightRequest := &base.GetLatestBlockRequest{}
	blocks := make([]*blockTypes.Block, 0)

	lastBlock, err := cc.BaseQuaryCli.GetLatestBlock(context.Background(), blockByHeightRequest)
	if err != nil {
		logger.Error("Failed to query the last block, err", err)
		return nil, nil, nil, err
	}
	logger.Info("Query block high success")

	for i := start; i < lastBlock.GetBlock().Header.Height; i++ {
		blockByHeightRequest := &base.GetBlockByHeightRequest{
			Height: i,
		}
		blockInfo, err := cc.BaseQuaryCli.GetBlockByHeight(context.Background(), blockByHeightRequest)
		if err != nil {
			logger.Error("Failed to query block data by block height, err:", err)
			return nil, nil, nil, err
		}

		block := blockInfo.GetBlock()
		blocks = append(blocks, block)
	}
	blocks = append(blocks, lastBlock.GetBlock())

	proposals := make(map[string]map[int64]struct{}, 0)
	signs := make(map[string]map[int64]struct{}, 0)
	addrMap := make(map[string]string)
	monikerMap := make(map[string]string)

	proposalAssignments := make([]*types.ProposalAssignment, 0)
	valSign := make([]*types.ValSign, 0)
	valSignMissed := make([]*types.ValSignMissed, 0)

	for _, monitorObj := range monitorObjs {
		proposals[strings.ToLower(monitorObj.OperatorAddrHex)] = make(map[int64]struct{}, 0)
		signs[strings.ToLower(monitorObj.OperatorAddrHex)] = make(map[int64]struct{}, 0)
		addrMap[strings.ToLower(monitorObj.OperatorAddrHex)] = monitorObj.OperatorAddr
		monikerMap[strings.ToLower(monitorObj.OperatorAddrHex)] = monitorObj.Moniker
	}

	for _, block := range blocks {
		if _, ok := addrMap[strings.ToLower(hex.EncodeToString(block.Header.GetProposerAddress()))]; ok {
			proposalAssignments = append(proposalAssignments, &types.ProposalAssignment{
				BlockHeight:  block.GetHeader().GetHeight(),
				Moniker:      monikerMap[strings.ToLower(hex.EncodeToString(block.Header.GetProposerAddress()))],
				OperatorAddr: addrMap[strings.ToLower(hex.EncodeToString(block.Header.GetProposerAddress()))],
				ChildTable:   block.GetHeader().GetHeight() % 10,
			})
		}

		for _, sign := range block.LastCommit.Signatures {
			if _, ok := addrMap[strings.ToLower(hex.EncodeToString(sign.ValidatorAddress))]; ok {
				signs[strings.ToLower(hex.EncodeToString(sign.ValidatorAddress))][block.GetHeader().Height] = struct{}{}

				valSign = append(valSign, &types.ValSign{
					BlockHeight:  block.GetHeader().GetHeight(),
					Moniker:      monikerMap[strings.ToLower(hex.EncodeToString(sign.ValidatorAddress))],
					OperatorAddr: addrMap[strings.ToLower(hex.EncodeToString(sign.ValidatorAddress))],
					Status:       1,
					DoubleSign:   false,
					ChildTable:   block.GetHeader().GetHeight() % 10,
				})
			}

		}
	}

	for _, monitorObj := range monitorObjs {
		for height := start; height <= lastBlock.GetBlock().Header.Height; height++ {
			if _, ok := signs[strings.ToLower(monitorObj.OperatorAddrHex)][height]; !ok {
				valSignMissed = append(valSignMissed, &types.ValSignMissed{
					Moniker:      monikerMap[strings.ToLower(monitorObj.OperatorAddrHex)],
					BlockHeight:  int(height),
					OperatorAddr: addrMap[strings.ToLower(monitorObj.OperatorAddrHex)],
				})
			}
		}
	}

	return proposalAssignments, valSign, valSignMissed, nil
}

func (cc *ChainCli) GetValRanking(monitorObjs []*types.MonitorObj, project string) ([]*types.ValRanking, error) {
	cons2Moniker := make(map[string]string, 0)
	cons2Operator := make(map[string]string, 0)
	for _, monitorObj := range monitorObjs {
		consAddr := utils.Operator2Cons(monitorObj.OperatorAddrHex, project)
		cons2Moniker[consAddr] = monitorObj.Moniker
		cons2Operator[consAddr] = monitorObj.OperatorAddr
	}
	queryValRankingRequest := &base.GetLatestValidatorSetRequest{}
	valsRanking, err := cc.BaseQuaryCli.GetLatestValidatorSet(context.Background(), queryValRankingRequest)
	if err != nil {
		logger.Error("Failed to query LatestValidatorSet, err:", err)
		return nil, err
	}
	valSet := make(map[string]*types.ValSet)
	for ranking, val := range valsRanking.GetValidators() {
		valSet[val.GetAddress()] = &types.ValSet{
			Validators: &base.Validator{
				Address:          val.GetAddress(),
				PubKey:           val.GetPubKey(),
				VotingPower:      val.GetVotingPower(),
				ProposerPriority: val.GetProposerPriority(),
			},
			Ranking: ranking + 1,
		}
	}

	valRanking := make([]*types.ValRanking, 0)
	for consAddr, moniker := range cons2Moniker {
		if val, ok := valSet[consAddr]; ok {
			valRanking = append(valRanking, &types.ValRanking{
				ChainName:       project,
				BlockHeight:     valsRanking.BlockHeight,
				Moniker:         moniker,
				OperatorAddr:    cons2Operator[consAddr],
				RealVotingPower: val.Validators.GetVotingPower(),
				Ranking:         val.Ranking,
			})
		} else {
			valRanking = append(valRanking, &types.ValRanking{
				ChainName:       project,
				Moniker:         moniker,
				OperatorAddr:    cons2Operator[consAddr],
				BlockHeight:     valsRanking.BlockHeight,
				RealVotingPower: 0,
				Ranking:         1000,
			})
		}
	}
	return valRanking, err
}

var logger = log.RPCLogger.WithField("module", "rpc")
