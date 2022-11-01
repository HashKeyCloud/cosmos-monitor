package injective

import (
	"fmt"
	"testing"
)

func TestGetValInfo(t *testing.T) {
	cc, err := NewInjectiveRpcCli("injective-grpc.polkachu.com:14390")
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
