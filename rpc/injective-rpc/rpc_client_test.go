package injective_rpc

import (
	"cosmosmonitor/types"
	"fmt"
	"testing"
)

func TestGetValInfo(t *testing.T) {
	cc, err := NewInjectiveRpcCli()
	if err != nil {
		logger.Error("err:", err)
	}
	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "injvaloper1g4d6dmvnpg7w7yugy6kplndp7jpfmf3krtschp")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetProposal(t *testing.T) {
	cc, err := NewInjectiveRpcCli()
	if err != nil {
		logger.Error("err:", err)
	}

	monitorObjs := make([]*types.MonitorObj, 0)
	monitorObjs = append(monitorObjs, &types.MonitorObj{
		Moniker:         "Figment",
		OperatorAddr:    "injvaloper1g4d6dmvnpg7w7yugy6kplndp7jpfmf3krtschp",
		OperatorAddrHex: "6087607e1e56f6ee7934abaf65834c92d618104c",
		SelfStakeAddr:   "inj1g4d6dmvnpg7w7yugy6kplndp7jpfmf3k5d9ak9",
	})
	monitors, _ := cc.GetProposal(monitorObjs)
	for _, monitor := range monitors {
		fmt.Println("proposal:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	cc, err := NewInjectiveRpcCli()
	if err != nil {
		logger.Error("err:", err)
	}

	monitorObjs := make([]*types.MonitorObj, 0)
	monitorObjs = append(monitorObjs, &types.MonitorObj{
		Moniker:         "Figment",
		OperatorAddr:    "injvaloper1g4d6dmvnpg7w7yugy6kplndp7jpfmf3krtschp",
		OperatorAddrHex: "6087607e1e56f6ee7934abaf65834c92d618104c",
		SelfStakeAddr:   "inj1g4d6dmvnpg7w7yugy6kplndp7jpfmf3k5d9ak9",
	})
	proposalAssignment, valSign, valSignMissed, _ := cc.GetValPerformance(18186100, monitorObjs)
	for _, monitor := range proposalAssignment {
		fmt.Println("proposalAssignment:", monitor)
	}

	for _, monitor := range valSign {
		fmt.Println("valSign:", monitor)
	}

	for _, monitor := range valSignMissed {
		fmt.Println("valSignMissed:", monitor)
	}
}
