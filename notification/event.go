package notification

import (
	"fmt"

	"github.com/spf13/viper"

	"cosmosmonitor/types"
)

type Event interface {
	Name() string
	Message() string
	IsEmpty() bool
}

type exception struct {
	validators []*struct {
		blockHeight int64
		moniker     string
	}
}

type ProposalException struct {
	proposals []*types.Proposal
}

type (
	ValJailedException   exception
	ValisActiveException exception
	SyncException        exception
)

func (e *ValJailedException) Name() string {
	return "Validator Jailed Exception\n"
}
func (e *ValJailedException) Message() string {
	var msg string
	if !e.IsEmpty() {
		msg = e.Name()
		for _, val := range e.validators {
			msg += fmt.Sprintf("%s validator has been jailed\n", val.moniker)
		}
	}
	return msg
}
func (e *ValJailedException) IsEmpty() bool {
	return 0 == len(e.validators)
}

func (e *ValisActiveException) Name() string {
	return "Validator InActive Exception\n"
}
func (e *ValisActiveException) Message() string {
	var msg string
	if !e.IsEmpty() {
		msg = e.Name()
		for _, val := range e.validators {
			msg += fmt.Sprintf("%s validator is Inactive\n", val.moniker)
		}
	}
	return msg
}
func (e *ValisActiveException) IsEmpty() bool {
	return 0 == len(e.validators)
}

func (e *SyncException) Name() string {
	return "Sync Exception \n"
}
func (e *SyncException) Message() string {
	proportion := viper.GetFloat64("alert.proportion")
	var msg string
	if !e.IsEmpty() {
		msg = e.Name()
		for _, val := range e.validators {
			msg += fmt.Sprintf("The %s validator has not signed for 5 consecutive blocks or the last 100 blocks without signature rate reaches %f at block height of %d. \n",
				val.moniker, proportion, val.blockHeight)
		}
	}
	return msg
}

func (e *SyncException) IsEmpty() bool {
	return 0 == len(e.validators)
}

func (e *ProposalException) Name() string {
	return "proposal Exception \n"
}

func (e *ProposalException) Message() string {
	var msg string
	if !e.IsEmpty() {
		msg = e.Name()
		for _, proposal := range e.proposals {
			msg += fmt.Sprintf("There is a new proposal\nThe proposal id is: %d \nThe voting start time is: %s \nThe acceptance time is: %s \nThe proposal content is: %s \n\n\n",
				proposal.ProposalId, proposal.VotingStartTime, proposal.VotingEndTime, proposal.Description)
		}
	}
	return msg
}

func (e *ProposalException) IsEmpty() bool {
	return 0 == len(e.proposals)
}

func ParseValJailedException(valJaileds []string) *ValJailedException {
	if len(valJaileds) == 0 {
		logger.Error("validator jailed is empty, please check")
		return nil
	}

	valJailedException := &ValJailedException{
		validators: make([]*struct {
			blockHeight int64
			moniker     string
		}, 0),
	}

	for _, valJailed := range valJaileds {
		valJailedException.validators = append(valJailedException.validators, &struct {
			blockHeight int64
			moniker     string
		}{moniker: valJailed})
	}
	return valJailedException
}

func ParseValisActiveException(valisActive []*types.ValIsActive) *ValisActiveException {
	if len(valisActive) == 0 {
		logger.Error("validator inActive is empty, please check")
		return nil
	}

	valisActiveException := &ValisActiveException{
		validators: make([]*struct {
			blockHeight int64
			moniker     string
		}, 0),
	}

	for _, valJailed := range valisActive {
		valisActiveException.validators = append(valisActiveException.validators, &struct {
			blockHeight int64
			moniker     string
		}{moniker: valJailed.Moniker})
	}
	return valisActiveException
}

func ParseProposalException(proposals []*types.Proposal) *ProposalException {
	if len(proposals) == 0 {
		logger.Error("proposal is empty, please check")
		return nil
	}
	proposalException := &ProposalException{
		proposals: make([]*types.Proposal, 0),
	}
	for _, proposal := range proposals {

		proposalException.proposals = append(proposalException.proposals, &types.Proposal{
			ProposalId:      proposal.ProposalId,
			VotingStartTime: proposal.VotingStartTime,
			VotingEndTime:   proposal.VotingEndTime,
			Description:     proposal.Description,
		})
	}
	return proposalException
}

func ParseSyncException(missedSign []*types.ValSignMissed) *SyncException {
	if len(missedSign) == 0 {
		logger.Error("validator missed sign is empty, please check")
		return nil
	}

	syncException := &SyncException{
		validators: make([]*struct {
			blockHeight int64
			moniker     string
		}, 0),
	}
	for _, valMissedSign := range missedSign {
		syncException.validators = append(syncException.validators, &struct {
			blockHeight int64
			moniker     string
		}{moniker: valMissedSign.OperatorAddr,
			blockHeight: int64(valMissedSign.BlockHeight)})

	}
	return syncException
}
