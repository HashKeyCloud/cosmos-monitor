package notification

import (
	"fmt"
	"testing"

	"cosmosmonitor/types"
)

func TestValJailedException(t *testing.T) {
	val := struct {
		blockHeight int64
		moniker     string
	}{blockHeight: 123,
		moniker: "monikerName"}
	vals := make([]*struct {
		blockHeight int64
		moniker     string
	}, 0)
	vals = append(vals, &val)
	ex := exception{
		vals,
	}
	vj := ParseValJailedException([]string{ex.validators[0].moniker})
	fmt.Println(vj.Message())
}

func TestParseValisActiveException(t *testing.T) {
	valIsActive := make([]*types.ValIsActive, 0)
	valIActive := &types.ValIsActive{
		"monikerName",
		1,
	}
	valIsActive = append(valIsActive, valIActive)
	vj := ParseValisActiveException(valIsActive)
	fmt.Println(vj.Message())
}

func TestParseSyncException(t *testing.T) {
	valSignMissed := make([]*types.ValSignMissed, 0)
	valSignMissed = append(valSignMissed, &types.ValSignMissed{
		"monikerName",
		"1",
		1,
	})
	vj := ParseSyncException(valSignMissed)
	fmt.Println(vj.Message())
}
