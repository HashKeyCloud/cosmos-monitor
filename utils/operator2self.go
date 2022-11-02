package utils

import (
	"fmt"
	"github.com/cosmos/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosBech32 "github.com/cosmos/cosmos-sdk/types/bech32"
)

func Operator2SelfAddr(operator string) string {
	hrp, bz, err := cosmosBech32.DecodeAndConvert(operator)
	if err != nil {
		fmt.Println("err:", err)
	}

	operatorBytes := []byte(sdk.ValAddress(bz))
	operatorByte, err := bech32.ConvertBits(operatorBytes, 8, 5, true)
	if err != nil {
		logger.Info("Failed to convert string to bytes, Error:", err)
	}

	var selfaddr string
	switch hrp {
	case "cosmosvaloper":
		selfaddr, err = bech32.Encode("cosmos", operatorByte) //cosmos
	case "evmosvaloper":
		selfaddr, err = bech32.Encode("evmos", operatorByte) // evmos
	case "injvaloper":
		selfaddr, err = bech32.Encode("inj", operatorByte) // injective
	case "nvaloper":
		selfaddr, err = bech32.Encode("n1", operatorByte) // nyx
	case "persistencevaloper":
		selfaddr, err = bech32.Encode("persistence", operatorByte) // persistence
	case "rizonvaloper":
		selfaddr, err = bech32.Encode("rizon", operatorByte) // rizon
	case "secretvaloper":
		selfaddr, err = bech32.Encode("secret", operatorByte) // secret
	case "sommvaloper":
		selfaddr, err = bech32.Encode("somm", operatorByte) // sommelier
	case "torivaloper":
		selfaddr, err = bech32.Encode("tori", operatorByte) // Teritori
	}

	if err != nil {
		logger.Error("Failed to convert hex format to bech32, err:", err)
	}

	return selfaddr
}
