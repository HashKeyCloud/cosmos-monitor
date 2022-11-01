package cosmos

import (
	"cosmosmonitor/types"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"testing"
	"time"
)

func TestGovInfo(t *testing.T) {
	cc, err := NewCosmosRpcCli("xxxx")
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
	cc, err := NewCosmosRpcCli("cosmos-grpc.polkachu.com:14990")
	if err != nil {
		logger.Error("err:", err)
	}
	monitorObj := make([]string, 0)
	monitorObj = append(monitorObj, "cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en")
	monitors, _ := cc.GetValInfo(monitorObj)
	for _, monitor := range monitors {
		fmt.Println("monitor:", monitor)
	}
}

func TestGetValPerformance(t *testing.T) {
	cc, err := NewCosmosRpcCli("xxxx")
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

func TestA(t *testing.T) {
	/*data := []byte("cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en")
	// Convert test data to base32:
	conv, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		fmt.Println("Error:", err)
	}
	encoded, err := bech32.Encode("cosmosvaloper", conv)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Show the encoded data.
	fmt.Println("Encoded Data:", encoded)*/

	//encoded := "cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en"
	//hrp, decoded, err := bech32.Decode(encoded)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//}
	//
	//fmt.Println("hrp", hrp)
	//fmt.Println("decoded:", hex.EncodeToString(decoded))

	hrp, bz, err := bech32.DecodeAndConvert("cosmosvaloper1c4k24jzduc365kywrsvf5ujz4ya6mwympnc4en")
	if err != nil {
		fmt.Println("err:", err)
	}

	if hrp != "cosmosvaloper" {
		fmt.Println("err:", err)
	}

	fmt.Println(hex.EncodeToString(bz))
}
