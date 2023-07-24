package utils

import (
	"encoding/hex"
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
		logger.Error("Failed to convert string to bytes, Error:", err)
	}

	var selfAddr string
	switch hrp {
	case "acrevaloper":
		selfAddr, err = bech32.Encode("acre", operatorByte) // acrechain
	case "akashvaloper":
		selfAddr, err = bech32.Encode("akash", operatorByte) // akash
	case "axelarvaloper":
		selfAddr, err = bech32.Encode("axelar", operatorByte) // axelar
	case "bandvaloper":
		selfAddr, err = bech32.Encode("band", operatorByte) //band
	case "cosmosvaloper":
		selfAddr, err = bech32.Encode("cosmos", operatorByte) // cosmos and cosmos's consumer chain
	case "evmosvaloper":
		selfAddr, err = bech32.Encode("evmos", operatorByte) // evmos
	case "injvaloper":
		selfAddr, err = bech32.Encode("inj", operatorByte) // injective
	case "junovaloper":
		selfAddr, err = bech32.Encode("juno", operatorByte) // juno
	case "neutronvaloper":
		selfAddr, err = bech32.Encode("neutron", operatorByte) // neutron
	case "nvaloper":
		selfAddr, err = bech32.Encode("n", operatorByte) // nyx
	case "okp4valoper":
		selfAddr, err = bech32.Encode("okp4", operatorByte) // okp4
	case "persistencevaloper":
		selfAddr, err = bech32.Encode("persistence", operatorByte) // persistence
	case "rizonvaloper":
		selfAddr, err = bech32.Encode("rizon", operatorByte) // rizon
	case "secretvaloper":
		selfAddr, err = bech32.Encode("secret", operatorByte) // secret
	case "sommvaloper":
		selfAddr, err = bech32.Encode("somm", operatorByte) // sommelier
	case "torivaloper":
		selfAddr, err = bech32.Encode("tori", operatorByte) // teritori
	case "xplavaloper":
		selfAddr, err = bech32.Encode("xpla", operatorByte) // xpla
	case "zetavaloper":
		selfAddr, err = bech32.Encode("zeta", operatorByte) // zeta
	}

	if err != nil {
		logger.Error("Failed to convert hex format to bech32, err:", err)
	}

	return selfAddr
}

func Operator2Cons(operatorHex, project string) string {
	data, err := hex.DecodeString(operatorHex)
	if err != nil {
		logger.Error("Converting Hex to Byte fails, err:", err)
	}
	// Convert test data to base32:
	conv, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		logger.Error("Converting string to bech32 fails, err:", err)
	}
	var consAddr string
	switch project {
	case "acrechain":
		consAddr, err = bech32.Encode("acrevalcons", conv) // acrechain
	case "akash":
		consAddr, err = bech32.Encode("akashvalcons", conv) // akash
	case "apollo":
		consAddr, err = bech32.Encode("cosmosvalcons", conv) // apollo
	case "axelar":
		consAddr, err = bech32.Encode("axelarvalcons", conv) // axelar
	case "band":
		consAddr, err = bech32.Encode("bandvalcons", conv) // band
	case "cosmos":
		consAddr, err = bech32.Encode("cosmosvalcons", conv) // cosmos
	case "evmos":
		consAddr, err = bech32.Encode("evmosvalcons", conv) // evmos
	case "gopher":
		consAddr, err = bech32.Encode("cosmosvalcons", conv) // gopher
	case "hero":
		consAddr, err = bech32.Encode("cosmosvalcons", conv) // hero
	case "injective":
		consAddr, err = bech32.Encode("injvalcons", conv) // injective
	case "juno":
		consAddr, err = bech32.Encode("junovalcons", conv) // juno
	case "neutronconsumer":
		consAddr, err = bech32.Encode("cosmosvalcons", conv) // neutron consumer chain
	case "neutron":
		consAddr, err = bech32.Encode("neutronvalcons", conv) // neutron
	case "nyx":
		consAddr, err = bech32.Encode("nvalcons", conv) // nyx
	case "okp4":
		consAddr, err = bech32.Encode("okp4valcons", conv) // okp4
	case "persistence":
		consAddr, err = bech32.Encode("persistencevalcons", conv) // persistence
	case "provider":
		consAddr, err = bech32.Encode("cosmosvalcons", conv) // provider
	case "rizon":
		consAddr, err = bech32.Encode("rizonvalcons", conv) // rizon
	case "secret":
		consAddr, err = bech32.Encode("secretvalcons", conv) // secret
	case "sommelier":
		consAddr, err = bech32.Encode("sommvalcons", conv) // sommelier
	case "sputnik":
		consAddr, err = bech32.Encode("cosmosvalcons", conv) // sputnik
	case "teritori":
		consAddr, err = bech32.Encode("torivalcons", conv) // teritori
	case "xpla":
		consAddr, err = bech32.Encode("xplavalcons", conv) // xpla
	case "zeta":
		consAddr, err = bech32.Encode("zetavalcons", conv) // xpla
	}

	if err != nil {
		logger.Error("Conversion to consAddr failed, err:", err)
	}
	return consAddr
}
