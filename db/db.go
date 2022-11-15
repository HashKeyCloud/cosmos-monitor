package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type DBCli interface {
	SaveValInfo(valInfo []*types.ValInfo) error
	BatchSaveValSign(valSigns []*types.ValSign) error
	BatchSaveValSignMissed(valSignMissed []*types.ValSignMissed) error
	BatchSaveProposalAssignments(proposalAssignments []*types.ProposalAssignment) error
	BatchSaveProposals(proposals []*types.Proposal) error
	BatchSaveSignNum(startBlock, endBlock int64, operatorAddrs []string) error
	BatchSaveUptime(startBlock, endBlock int64, operatorAddrs []string) error
	BatchSaveMissedSignNum(startBlock, endBlock int64, operatorAddrs []string) error
	BatchSaveProposalsNum(startBlock, endBlock int64, operatorAddrs []string) error
	GetBlockHeightFromDb(project string) (int64, error)
	GetValSignMissedFromDb(start, end int64) ([]*types.ValSignMissed, error)
	GetValMoniker() ([]*types.ValMoniker, error)
	GetMonitorObj() ([]*types.MonitorObj, error)
	BatchSaveValStats(start, end int64) error
	BatchSaveRanking(valRankings []*types.ValRanking) error
}

type DbCli struct {
	Conn *sqlx.DB
}

func InitDB(dbconf *types.DatabaseConfig) (*DbCli, error) {
	port, _ := strconv.Atoi(dbconf.Port)
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbconf.Host, port, dbconf.Username, dbconf.Password, dbconf.Name)
	db, err := sqlx.Connect("postgres", pdqlInfo)
	if err != nil {
		logger.Errorf("Connected failed.err:%v\n", err)
		return nil, err
	}

	dbConnectionTimeout := time.NewTimer(15 * time.Second)
	go func() {
		<-dbConnectionTimeout.C
		logger.Fatalf("timeout while connecting to the database")
	}()
	err = db.Ping()
	if err != nil {
		logger.Errorf("ping db fail, err:%v", err)
	}

	dbConnectionTimeout.Stop()

	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Second * 60)
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(200)

	logger.Info("Successfully connected!")
	return &DbCli{
		Conn: db,
	}, nil
}

func (dc *DbCli) SaveValInfo(valInfo []*types.ValInfo) error {
	logger.Info("begin save validator info")
	batchSize := 500
	for b := 0; b < len(valInfo); b += batchSize {
		logger.Infof("Start saving %d batch of validator info\n", b+1)
		start := b
		end := b + batchSize
		if len(valInfo) < end {
			end = len(valInfo)
		}
		numArgs := 16
		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*numArgs)

		for i, v := range valInfo[start:end] {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*numArgs+1, i*numArgs+2, i*numArgs+3, i*numArgs+4, i*numArgs+5, i*numArgs+6, i*numArgs+7, i*numArgs+8, i*numArgs+9,
				i*numArgs+10, i*numArgs+11, i*numArgs+12, i*numArgs+13, i*numArgs+14, i*numArgs+15, i*numArgs+16))
			valueArgs = append(valueArgs, v.Moniker)
			valueArgs = append(valueArgs, v.OperatorAddr)
			valueArgs = append(valueArgs, v.OperatorAddrHex)
			valueArgs = append(valueArgs, v.SelfStakeAddr)
			valueArgs = append(valueArgs, v.RewardAddr)
			valueArgs = append(valueArgs, v.Jailed)
			valueArgs = append(valueArgs, v.Status)
			valueArgs = append(valueArgs, v.VotingPower)
			valueArgs = append(valueArgs, v.Identity)
			valueArgs = append(valueArgs, v.Website)
			valueArgs = append(valueArgs, v.Details)
			valueArgs = append(valueArgs, v.SecurityContact)
			valueArgs = append(valueArgs, v.CommissionRates)
			valueArgs = append(valueArgs, v.MaxRate)
			valueArgs = append(valueArgs, v.MaxChangeRate)
			valueArgs = append(valueArgs, v.MinSelfDelegation)
		}

		sql := fmt.Sprintf(`
			INSERT INTO val_info (
				moniker,
				operator_addr,
				operator_addr_hex,
				self_stake_addr,
				reward_addr,
				jailed,
				status,
				voting_power,
				identity,
				website,
				details,
				security_contact,
				commission_rates,
				max_rate,
				max_change_rate,
				min_self_delegation
			) 
			VALUES %v
			ON  CONFLICT (operator_addr) DO UPDATE SET
				moniker = EXCLUDED.moniker,
				operator_addr_hex = EXCLUDED.operator_addr_hex,
				self_stake_addr = EXCLUDED.self_stake_addr,
				reward_addr = EXCLUDED.reward_addr,
				jailed = EXCLUDED.jailed,
				status = EXCLUDED.status,
				voting_power = EXCLUDED.voting_power,
				identity = EXCLUDED.identity,
				website = EXCLUDED.website,
				details = EXCLUDED.details,
				security_contact = EXCLUDED.security_contact,
				commission_rates = EXCLUDED.commission_rates,
				max_rate = EXCLUDED.max_rate,
				max_change_rate = EXCLUDED.max_change_rate,
				min_self_delegation = EXCLUDED.min_self_delegation;
		`, strings.Join(valueStrings, ","))
		_, err := dc.Conn.Exec(sql, valueArgs...)
		if err != nil {
			logger.Errorf("saving validator batch %v fail. err:%v \n", b+1, err)
			return err
		}
		logger.Info("saving validator completed")
	}
	return nil
}

func (dc *DbCli) BatchSaveValSign(valSigns []*types.ValSign) error {
	batchSize := 500
	for b := 0; b < len(valSigns); b += batchSize {
		logger.Infof("Start saving %d batch of validator sign\n", b+1)
		start := b
		end := b + batchSize
		if len(valSigns) < end {
			end = len(valSigns)
		}
		numArgs := 6
		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*numArgs)

		for i, v := range valSigns[start:end] {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
				i*numArgs+1, i*numArgs+2, i*numArgs+3, i*numArgs+4, i*numArgs+5, i*numArgs+6))
			valueArgs = append(valueArgs, v.Moniker)
			valueArgs = append(valueArgs, v.OperatorAddr)
			valueArgs = append(valueArgs, v.BlockHeight)
			valueArgs = append(valueArgs, v.Status)
			valueArgs = append(valueArgs, v.DoubleSign)
			valueArgs = append(valueArgs, v.BlockHeight%10)
		}

		sql := fmt.Sprintf(`
			INSERT INTO val_sign_p (
				moniker,
				operator_addr,
				block_height,
				status,
				double_sign,
				child_table
			)
			VALUES %v
			ON  CONFLICT (operator_addr, block_height, child_table) DO UPDATE SET
				moniker = EXCLUDED.moniker,
				status = EXCLUDED.status,
				double_sign = EXCLUDED.double_sign;
		`, strings.Join(valueStrings, ","))
		_, err := dc.Conn.Exec(sql, valueArgs...)
		if err != nil {
			logger.Errorf("saving validator sign batch %v fail. err:%v \n", b+1, err)
			return err
		}

		logger.Infof("saving validator sign %v completed\n", b+1)
	}
	return nil
}
func (dc *DbCli) BatchSaveValSignMissed(valSignMissed []*types.ValSignMissed) error {
	batchSize := 500
	for b := 0; b < len(valSignMissed); b += batchSize {
		logger.Infof("Start saving %d batch of validator sign missed\n", b+1)
		start := b
		end := b + batchSize
		if len(valSignMissed) < end {
			end = len(valSignMissed)
		}
		numArgs := 3
		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*numArgs)

		for i, v := range valSignMissed[start:end] {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)",
				i*numArgs+1, i*numArgs+2, i*numArgs+3))
			valueArgs = append(valueArgs, v.Moniker)
			valueArgs = append(valueArgs, v.OperatorAddr)
			valueArgs = append(valueArgs, v.BlockHeight)
		}

		sql := fmt.Sprintf(`
			INSERT INTO val_sign_missed (
				moniker,
				operator_addr,
				block_height
			)
			VALUES %v
			ON  CONFLICT (operator_addr, block_height) DO UPDATE SET
				moniker = EXCLUDED.moniker;
		`, strings.Join(valueStrings, ","))
		_, err := dc.Conn.Exec(sql, valueArgs...)
		if err != nil {
			logger.Errorf("saving validator sign missed batch %v fail. err:%v \n", b+1, err)
			return err
		}

		logger.Infof("saving validator sign missed %v completed\n", b+1)
	}
	return nil
}
func (dc *DbCli) BatchSaveProposalAssignments(proposalAssignments []*types.ProposalAssignment) error {
	batchSize := 500
	for b := 0; b < len(proposalAssignments); b += batchSize {
		logger.Infof("Start saving %d batch of ProposalAssignments\n", b+1)
		start := b
		end := b + batchSize
		if len(proposalAssignments) < end {
			end = len(proposalAssignments)
		}
		numArgs := 4
		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*numArgs)

		for i, v := range proposalAssignments[start:end] {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)",
				i*numArgs+1, i*numArgs+2, i*numArgs+3, i*numArgs+4))
			valueArgs = append(valueArgs, v.Moniker)
			valueArgs = append(valueArgs, v.OperatorAddr)
			valueArgs = append(valueArgs, v.BlockHeight)
			valueArgs = append(valueArgs, v.BlockHeight%10)
		}

		sql := fmt.Sprintf(`
			INSERT INTO proposal_assignments_p (
				moniker,
				operator_addr,
				block_height,
				child_table
			)
			VALUES %v
			ON  CONFLICT (operator_addr, block_height, child_table) DO UPDATE SET
				moniker = EXCLUDED.moniker;
		`, strings.Join(valueStrings, ","))
		_, err := dc.Conn.Exec(sql, valueArgs...)
		if err != nil {
			logger.Errorf("saving ProposalAssignments %v fail. err:%v \n", b+1, err)
			return err
		}

		logger.Infof("saving ProposalAssignments missed %v completed\n", b+1)
	}
	return nil
}
func (c *DbCli) BatchSaveProposals(proposals []*types.Proposal) error {
	batchSize := 10
	for b := 0; b < len(proposals); b += batchSize {
		logger.Infof("Start saving %d batch of Proposals\n", b+1)
		start := b
		end := b + batchSize
		if len(proposals) < end {
			end = len(proposals)
		}
		numArgs := 7
		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*numArgs)

		for i, v := range proposals[start:end] {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*numArgs+1, i*numArgs+2, i*numArgs+3, i*numArgs+4, i*numArgs+5, i*numArgs+6, i*numArgs+7))
			valueArgs = append(valueArgs, v.ProposalId)
			valueArgs = append(valueArgs, v.VotingStartTime)
			valueArgs = append(valueArgs, v.VotingEndTime)
			valueArgs = append(valueArgs, v.Description)
			valueArgs = append(valueArgs, v.Moniker)
			valueArgs = append(valueArgs, v.OperatorAddr)
			valueArgs = append(valueArgs, v.Status)
		}

		sql := fmt.Sprintf(`
			INSERT INTO proposal (
				proposal_id,
				voting_start_time,
				voting_end_time,
				description,
				moniker,
				operator_addr,
				status
			)
			VALUES %v
			ON  CONFLICT (proposal_id, operator_addr) DO UPDATE SET
				voting_start_time = EXCLUDED.voting_start_time,
				voting_end_time = EXCLUDED.voting_end_time,
				description = EXCLUDED.description,
				moniker = EXCLUDED.moniker,
				status = EXCLUDED.status;
		`, strings.Join(valueStrings, ","))
		_, err := c.Conn.Exec(sql, valueArgs...)
		if err != nil {
			logger.Errorf("saving Proposals %v fail. err:%v \n", b+1, err)
			return err
		}

		logger.Infof("saving Proposals %v completed\n", b+1)
	}
	return nil
}

func (dc *DbCli) BatchSaveSignNum(startBlock, endBlock int64, operatorAddrs []string) error {
	for _, operatorAddr := range operatorAddrs {
		logger.Infof("Begin SaveSignNum for %v validator succeeded\n", operatorAddr)
		sql := `
			INSERT INTO val_stats (moniker, operator_addr, start_block, end_block, sign_num)
			(
				SELECT moniker, operator_addr, $2, $3, COUNT(*) FROM val_sign_p
				WHERE operator_addr = $1 AND block_height >= $2 AND block_height <= $3
				GROUP BY moniker, operator_addr
			)
			ON CONFLICT(operator_addr, start_block, end_block) do update set  moniker = excluded.moniker, sign_num = excluded.sign_num;
		`
		_, err := dc.Conn.Exec(sql, operatorAddr, startBlock, endBlock)
		if err != nil {
			logger.Errorf("Failed to save SaveSignNum for %v validator, err: %v \n", operatorAddr, err)
			return err
		}

		logger.Infof("Save SaveSignNum for %v validator succeeded\n", operatorAddr)
	}

	return nil
}

func (dc *DbCli) BatchSaveUptime(startBlock, endBlock int64, operatorAddrs []string) error {
	logger.Infof("begin save uptime")
	valSignNum := make([]*types.ValSignNum, 0)
	valSignNumMap := make(map[string]float64)
	sql := `SELECT operator_addr, sign_num FROM val_stats WHERE start_block = $1 AND end_block = $2`
	err := dc.Conn.Select(&valSignNum, sql, startBlock, endBlock)
	if err != nil {
		logger.Errorf("Failed to query validator missed attestation num, err:%v\n", err)
		return err
	}
	for _, v := range valSignNum {
		valSignNumMap[v.OperatorAddr] = float64(v.SignNum)
	}

	batchSize := 10
	for b := 0; b < len(operatorAddrs); b += batchSize {
		logger.Infof("Start saving %d batch of Proposals\n", b+1)
		start := b
		end := b + batchSize
		if len(operatorAddrs) < end {
			end = len(operatorAddrs)
		}
		numArgs := 4
		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*numArgs)

		for i, v := range operatorAddrs[start:end] {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)",
				i*numArgs+1, i*numArgs+2, i*numArgs+3, i*numArgs+4))
			valueArgs = append(valueArgs, v)
			valueArgs = append(valueArgs, startBlock)
			valueArgs = append(valueArgs, endBlock)
			valueArgs = append(valueArgs, valSignNumMap[v]/(float64(endBlock)-float64(startBlock)+1))
		}

		sql := fmt.Sprintf(`
			INSERT INTO val_stats (
				operator_addr,
				start_block,
				end_block,
				uptime
			)
			VALUES %v
			ON  CONFLICT (operator_addr, start_block, end_block) DO UPDATE SET
				uptime = EXCLUDED.uptime;
		`, strings.Join(valueStrings, ","))
		_, err := dc.Conn.Exec(sql, valueArgs...)
		if err != nil {
			logger.Errorf("saving uptime %v fail. err:%v \n", b+1, err)
			return err
		}
		logger.Infof("saving uptime %v completed\n", b+1)
	}
	return nil
}

func (dc *DbCli) BatchSaveMissedSignNum(startBlock, endBlock int64, operatorAddrs []string) error {
	for _, operatorAddr := range operatorAddrs {
		logger.Infof("Begin MissedSignNum for %v validator succeeded\n", operatorAddr)
		sql := `
			INSERT INTO val_stats (moniker, operator_addr, start_block, end_block, missed_sign_num)
			(
				SELECT moniker, operator_addr, $2, $3, COUNT(*) FROM val_sign_missed
				WHERE operator_addr = $1 AND block_height >= $2 AND block_height <= $3
				GROUP BY moniker, operator_addr
			)
			ON CONFLICT(operator_addr, start_block, end_block) do update set  moniker = excluded.moniker, missed_sign_num = excluded.missed_sign_num;
		`
		_, err := dc.Conn.Exec(sql, operatorAddr, startBlock, endBlock)
		if err != nil {
			logger.Errorf("Failed to save MissedSignNum for %v validator, err: %v \n", operatorAddr, err)
			return err
		}

		logger.Infof("Save MissedSignNum for %v validator succeeded\n", operatorAddr)
	}

	return nil
}

func (dc *DbCli) BatchSaveProposalsNum(startBlock, endBlock int64, operatorAddrs []string) error {
	for _, operatorAddr := range operatorAddrs {
		logger.Infof("Begin ProposalsNum for %v validator succeeded\n", operatorAddr)
		sql := `
			INSERT INTO val_stats (moniker, operator_addr, start_block, end_block, proposals_num)
			(
				SELECT moniker, operator_addr, $2, $3, COUNT(*) FROM proposal_assignments_p
				WHERE operator_addr = $1 AND block_height >= $2 AND block_height <= $3
				GROUP BY moniker, operator_addr
			)
			ON CONFLICT(operator_addr, start_block, end_block) do update set  moniker = excluded.moniker, proposals_num = excluded.proposals_num;
		`
		_, err := dc.Conn.Exec(sql, operatorAddr, startBlock, endBlock)
		if err != nil {
			logger.Errorf("Failed to save ProposalsNum for %v validator, err: %v \n", operatorAddr, err)
			return err
		}

		logger.Infof("Save ProposalsNum for %v validator succeeded\n", operatorAddr)
	}

	return nil
}

func (dc *DbCli) GetBlockHeightFromDb(project string) (int64, error) {
	var minHeight int64
	dbHeight := make([]types.MaxBlockHeight, 0)
	sqld := `select (select max(block_height) from val_sign_p) max_block_height_sign,
       max(block_height) max_block_height_missed from val_sign_p;`
	err := dc.Conn.Select(&dbHeight, sqld)
	if err != nil {
		logger.Errorf("Failed to query block height from db, err:%v\n", err)
		return 0, err
	}

	for _, height := range dbHeight {
		var blockHeight int64
		if height.MaxBlockHeightSign.Int64 >= height.MaxBlockHeightMissed.Int64 {
			blockHeight = height.MaxBlockHeightSign.Int64
		} else {
			blockHeight = height.MaxBlockHeightMissed.Int64
		}
		if minHeight == 0 {
			minHeight = blockHeight
		} else {
			if minHeight > blockHeight {
				minHeight = blockHeight
			}
		}
	}
	if minHeight == 0 {
		startingBlockHeight := fmt.Sprintf("alert.%sStartingBlockHeight", project)
		minHeight = int64(viper.GetInt(startingBlockHeight))
	}
	return minHeight, nil
}

func (dc *DbCli) GetValSignMissedFromDb(start, end int64) ([]*types.ValSignMissed, error) {
	valSignMissed := make([]*types.ValSignMissed, 0)
	sqld := `SELECT block_height, operator_addr FROM val_sign_missed WHERE block_height >= $1 AND block_height <= $2;`
	err := dc.Conn.Select(&valSignMissed, sqld, start, end)
	if err != nil {
		logger.Errorf("Failed to query validator sign missed, err:%v\n", err)
		return nil, err
	}
	logger.Info("query validator sign missed successful")
	return valSignMissed, nil
}

func (dc *DbCli) GetValMoniker() ([]*types.ValMoniker, error) {
	valsMoniker := make([]*types.ValMoniker, 0)
	sqld := `SELECT moniker, operator_addr FROM val_info;`
	err := dc.Conn.Select(&valsMoniker, sqld)
	if err != nil {
		logger.Errorf("Failed to query validator sign missed, err:%v\n", err)
		return nil, err
	}
	logger.Info("query validator sign missed successful")
	return valsMoniker, nil
}

func (dc *DbCli) BatchSaveValStats(start, end int64) error {
	if start < 0 && end > 0 {
		start = 0
	} else if start < 0 && end == 0 {
		return errors.New("The starting block height is negative and the ending block height is 0")
	}
	allVal := make([]string, 0)
	sqld := `SELECT operator_addr FROM val_info`
	err := dc.Conn.Select(&allVal, sqld)
	if err != nil {
		logger.Error("Failed to query all validator, err:", err)
		return err
	}
	err = dc.BatchSaveSignNum(start, end, allVal)
	if err != nil {
		logger.Errorf("When the block height is %d to %d, saving the number of validator signatures failed, err:%v \n", start, end, err)
	}

	err = dc.BatchSaveMissedSignNum(start, end, allVal)
	if err != nil {
		logger.Errorf("When the block height is %d to %d, it fails to save the number of unsigned validators, err:%v \n", start, end, err)
	}

	err = dc.BatchSaveProposalsNum(start, end, allVal)
	if err != nil {
		logger.Errorf("When the block height is %d to %d, it fails to save the number of blocks produced by the validator, err:%v \n", start, end, err)
	}

	err = dc.BatchSaveUptime(start, end, allVal)
	if err != nil {
		logger.Errorf("Failed to save validator signature rate when block height is %d to %d, err:%v \n", start, end, err)
	}

	return nil
}

func (dc *DbCli) GetMonitorObj() ([]*types.MonitorObj, error) {
	monitorObjs := make([]*types.MonitorObj, 0)
	sqld := `SELECT moniker, operator_addr, operator_addr_hex, self_stake_addr FROM val_info;`
	err := dc.Conn.Select(&monitorObjs, sqld)
	if err != nil {
		logger.Errorf("Failed to query monitor object, err:%v\n", err)
		return nil, err
	}
	logger.Info("query monitor object successful")
	return monitorObjs, nil
}

func (dc *DbCli) BatchSaveRanking(valRankings []*types.ValRanking) error {
	logger.Info("begin save validator votingPower and ranking")
	batchSize := 500
	for b := 0; b < len(valRankings); b += batchSize {
		logger.Infof("Start saving %d batch of validator votingPower and ranking\n", b+1)
		start := b
		end := b + batchSize
		if len(valRankings) < end {
			end = len(valRankings)
		}
		numArgs := 5
		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*numArgs)

		for i, v := range valRankings[start:end] {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
				i*numArgs+1, i*numArgs+2, i*numArgs+3, i*numArgs+4, i*numArgs+5))
			valueArgs = append(valueArgs, v.Moniker)
			valueArgs = append(valueArgs, v.OperatorAddr)
			valueArgs = append(valueArgs, v.BlockHeight)
			valueArgs = append(valueArgs, v.RealVotingPower)
			valueArgs = append(valueArgs, v.Ranking)
		}

		sql := fmt.Sprintf(`
			INSERT INTO val_ranking (
				moniker,
				operator_addr,
				block_height,
				real_voting_power,
				ranking
			) 
			VALUES %v
			ON  CONFLICT (operator_addr, block_height) DO UPDATE SET
				moniker = EXCLUDED.moniker,
				real_voting_power = EXCLUDED.real_voting_power,
				ranking = EXCLUDED.ranking;
		`, strings.Join(valueStrings, ","))
		_, err := dc.Conn.Exec(sql, valueArgs...)
		if err != nil {
			logger.Errorf("saving validator votingPower and ranking batch %v fail. err:%v \n", b+1, err)
			return err
		}
		logger.Info("saving validator votingPower and ranking completed")
	}
	return nil
}

var logger = log.DBLogger.WithField("module", "db")
