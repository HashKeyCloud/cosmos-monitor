package rpc

import (
	"cosmosmonitor/types"
	"fmt"
	"testing"
	"time"
)

func TestGovInfo(t *testing.T) {
	cc, err := NewRpcCli("xxxx")
	if err != nil {
		logger.Error("err", err)
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	monitorObj := make([]*types.MonitorObj, 0)
	monitorObj = append(monitorObj, &types.MonitorObj{
		"monikerName",
		"cosmosvaloper1xxxx",
		"xxxx",
		"cosmosxxxx",
	})
	proposals, _ := cc.GetProposal(monitorObj)
	for _, p := range proposals {
		fmt.Println(p)
	}

}

func TestGetValInfo(t *testing.T) {
	cc, err := NewRpcCli("xxxx")
	if err != nil {
		logger.Error("err:", err)
	}
	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "cosmosvaloper1xxxx")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	cc, err := NewRpcCli("xxxx")
	if err != nil {
		logger.Error("err:", err)
	}
	monitorObj := make([]*types.MonitorObj, 0)
	m1 := &types.MonitorObj{
		"monikerName",
		"cosmosvaloper1xxxx",
		"xxxx",
		"cosmosxxxx",
	}
	monitorObj = append(monitorObj, m1)
	proposalAssignments, signs, signsmissed, _ := cc.GetValPerformance(12572343, monitorObj)
	for _, proposalAssignment := range proposalAssignments {
		fmt.Println("proposalAssignment:", proposalAssignment)
	}
	for _, sign := range signs {
		fmt.Println("sign:", sign)
	}

	for _, signmissed := range signsmissed {
		fmt.Println("sign missed:", signmissed)
	}
}
