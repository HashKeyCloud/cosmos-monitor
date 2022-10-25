package db

import (
	"cosmosmonitor/types"
	"fmt"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestBaveValInfo(t *testing.T) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"xxxxx", 5432, "ubuntu", "xxxxx", "cosmos_monitor")
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
	}
	dbCli := DbCli{db}

	valsInfo := make([]*types.ValInfo, 0)
	valInfo := &types.ValInfo{
		Moniker:           "monikerName",
		OperatorAddr:      "cosmosvaloperxxxx",
		OperatorAddrHex:   "xxxxxx",
		SelfStakeAddr:     "cosmos1xxxxx",
		RewardAddr:        "cosmosxxxxx",
		Jailed:            false,
		Status:            3,
		VotingPower:       "1000000000000",
		Identity:          "123-123",
		Website:           "www.website.com",
		Details:           "xxxxx",
		SecurityContact:   "xxx",
		CommissionRates:   0.1,
		MaxRate:           1.0,
		MaxChangeRate:     1.0,
		MinSelfDelegation: "1",
	}
	valsInfo = append(valsInfo, valInfo)
	err = dbCli.SaveValInfo(valsInfo)
	if err != nil {
		logger.Error("err：", err)
	}
}

func TestBatchSaveValSign(t *testing.T) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"xxxxx", 5432, "ubuntu", "xxxx", "cosmos_monitor")
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
	}
	dbCli := DbCli{db}
	valSigns := make([]*types.ValSign, 0)

	for i := 0; i < 10; i++ {
		valInfo := &types.ValSign{
			"monikerName",
			"cosmosvaloperxxxxx",
			int64(i),
			1,
			false,
			int64(i % 10),
		}
		valSigns = append(valSigns, valInfo)
	}

	err = dbCli.BatchSaveValSign(valSigns)
	if err != nil {
		logger.Error("err：", err)
	}
}

func TestBatchValSignMissed(t *testing.T) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"xxxx", 5432, "ubuntu", "xxxxx", "cosmos_monitor")
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
	}
	dbCli := DbCli{db}
	valSigns := make([]*types.ValSignMissed, 0)

	for i := 0; i < 10; i++ {
		valInfo := &types.ValSignMissed{
			"moniker",
			"cosmosvaloperxxxx",
			i,
		}
		valSigns = append(valSigns, valInfo)
	}

	err = dbCli.BatchSaveValSignMissed(valSigns)
	if err != nil {
		logger.Error("err：", err)
	}
}

//
func TestBatchProposalAssignments(t *testing.T) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"xxxx", 5432, "ubuntu", "xxx", "cosmos_monitor")
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
	}
	dbCli := DbCli{db}
	valSigns := make([]*types.ProposalAssignment, 0)

	for i := 0; i < 10; i++ {
		valInfo := &types.ProposalAssignment{
			"monikerName",
			"cosmosvaloperxxxxx",
			int64(i),
			1,
		}
		valSigns = append(valSigns, valInfo)
	}

	err = dbCli.BatchSaveProposalAssignments(valSigns)
	if err != nil {
		logger.Error("err：", err)
	}
}

func TestBatchProposals(t *testing.T) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"xxxx", 5432, "ubuntu", "xxxx", "cosmos_monitor")
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
	}
	dbCli := DbCli{db}
	valSigns := make([]*types.Proposal, 0)

	for i := 0; i < 10; i++ {
		valInfo := &types.Proposal{
			int64(i),
			"2022-12-12",
			"2022-12-24",
			"afduakfh",
			"monikerName",
			"cosmosxxxxx",
			1,
		}
		valSigns = append(valSigns, valInfo)
	}

	err = dbCli.BatchSaveProposals(valSigns)
	if err != nil {
		logger.Error("err：", err)
	}
}

func TestValStats(t *testing.T) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"xxxx", 5432, "ubuntu", "xxx", "cosmos_monitor")
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
	}
	dbCli := DbCli{db}
	val := []string{"cosmosvaloperxxxxx"}
	dbCli.BatchSaveMissedSignNum(0, 11, val)
	dbCli.BatchSaveSignNum(0, 11, val)
	dbCli.BatchSaveProposalsNum(0, 11, val)
	dbCli.BatchSaveUptime(0, 11, val)
}

func TestGetBlockHeightFromDb(t *testing.T) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"xxxx", 5432, "ubuntu", "xxx", "cosmos_monitor")
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
	}
	dbCli := DbCli{
		db,
	}
	height, err := dbCli.GetBlockHeightFromDb()
	if err != nil {
		logger.Error("get block height fail, err:", err)
	}
	logger.Info("height:", height)
}
