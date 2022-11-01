package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosBech32 "github.com/cosmos/cosmos-sdk/types/bech32"
)

func Operator2SelfAddr(operator string) string {
	operatorHex, err := sdk.ValAddressFromBech32(operator)
	fmt.Println("operatorHex:", operatorHex)
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

func Operator2SelfAddrInj(operator string) string {
	hrp, bz, err := cosmosBech32.DecodeAndConvert(operator)
	if err != nil {
		fmt.Println("err:", err)
	}

	if hrp != "injvaloper" {
		fmt.Println("err:", err)
	}

	operatorHex := hex.EncodeToString(bz)
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

	selfaddr, err := bech32.Encode("inj", operatorByte)
	if err != nil {
		logger.Error("Failed to convert hex format to bech32, err:", err)
	}

	return selfaddr
}
