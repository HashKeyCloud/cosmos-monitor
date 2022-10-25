package utils

import (
	"github.com/cosmos/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Operator2SelfAdde(operator string) string {
	operatorHex, err := sdk.ValAddressFromBech32(operator)
	if err != nil {
		logger.Error("Failed to convert bech32 format to Hex format, err:", err)
	}
	logger.Info("Bech32 format converted to Hex format successfully")

	logger.Info("operatorHex:", operatorHex)
	operatorBytes := []byte(operatorHex)

	operatorByte, err := bech32.ConvertBits(operatorBytes, 8, 5, true)
	if err != nil {
		logger.Info("Failed to convert string to bytes, Error:", err)
	}

	selfaddr, err := bech32.Encode("cosmos", operatorByte)
	if err != nil {
		logger.Error("Failed to convert hex format to bech32, err:", err)
	}

	return selfaddr
}
