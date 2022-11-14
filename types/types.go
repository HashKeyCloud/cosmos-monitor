package types

import "database/sql"

type DatabaseConfig struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
}

type ValInfo struct {
	Moniker           string  `db:"moniker"`
	OperatorAddr      string  `db:"operator_addr"`
	OperatorAddrHex   string  `db:"operator_addr_hex"`
	SelfStakeAddr     string  `db:"self_stake_addr"`
	RewardAddr        string  `db:"reward_addr"`
	Jailed            bool    `db:"jailed"`
	Status            int32   `db:"status"` //  1 = UNBONDED = inactive, 2 = UNBONDING = jail, 3 = BONDED = active
	VotingPower       string  `db:"voting_power"`
	Identity          string  `db:"identity"`
	Website           string  `db:"website"`
	Details           string  `db:"details"`
	SecurityContact   string  `db:"security_contact"`
	CommissionRates   float64 `db:"commission_rates"`
	MaxRate           float64 `db:"max_rate"`
	MaxChangeRate     float64 `db:"max_change_rate"`
	MinSelfDelegation string  `db:"min_self_delegation"`
}

type ValStats struct {
	Moniker       string  `db:"moniker"`
	OperatorAddr  string  `db:"operator_addr"`
	StartBlock    int     `db:"start_block"`
	EndBlock      int     `db:"end_block"`
	SignNum       int     `db:"sign_num"`
	MissedSignNum int     `db:"missed_sign_num"`
	ProposalsNum  int     `db:"proposals_num"`
	Uptime        float64 `db:"uptime"`
}

type ValSignNum struct {
	Moniker      string `db:"moniker"`
	OperatorAddr string `db:"operator_addr"`
	SignNum      int    `db:"sign_num"`
}

type ValMissedSignNum struct {
	Moniker       string `db:"moniker"`
	OperatorAddr  string `db:"operator_addr"`
	MissedSignNum int    `db:"missed_sign_num"`
}

type ValProposalsNum struct {
	Moniker      string `db:"moniker"`
	OperatorAddr string `db:"operator_addr"`
	ProposalsNum int    `db:"proposals_num"`
}

type ValSign struct {
	Moniker      string `db:"moniker"`
	OperatorAddr string `db:"operator_addr"`
	BlockHeight  int64  `db:"block_height"`
	Status       int    `db:"status"`
	DoubleSign   bool   `db:"double_sign"`
	ChildTable   int64  `db:"child_table"`
}

type ValSignMissed struct {
	ChainName    string
	Moniker      string `db:"moniker"`
	OperatorAddr string `db:"operator_addr"`
	BlockHeight  int    `db:"block_height"`
}

type ProposalAssignment struct {
	Moniker      string `db:"moniker"`
	OperatorAddr string `db:"operator_addr"`
	BlockHeight  int64  `db:"block_height"`
	ChildTable   int64  `db:"child_table"`
}

type Proposal struct {
	ChainName       string
	ProposalId      int64  `db:"proposal_id"`
	VotingStartTime string `db:"voting_start_time"`
	VotingEndTime   string `db:"voting_end_time"`
	Description     string `db:"description"`
	Moniker         string `db:"moniker"`
	OperatorAddr    string `db:"operator_addr"`
	Status          int32  `db:"status"`
}

type MonitorObj struct {
	Moniker         string `db:"moniker"`
	OperatorAddr    string `db:"operator_addr"`
	OperatorAddrHex string `db:"operator_addr_hex"`
	SelfStakeAddr   string `db:"self_stake_addr"`
}

type Height struct {
	OperatorAddr    string `db:"operator_addr"`
	OperatorAddrHex string `db:"operator_addr_hex"`
	BlockHeight     int64  `db:"block_height"`
}

type CaredData struct {
	ChainName           string
	ValInfos            []*ValInfo
	Proposals           []*Proposal
	ProposalAssignments []*ProposalAssignment
	ValSigns            []*ValSign
	ValSignMisseds      []*ValSignMissed
}

type ValIsActive struct {
	ChainName string
	Moniker   string `db:"moniker"`
	Status    int32  `db:"status"`
}

type ValIsJail struct {
	ChainName string
	Moniker   string `db:"moniker"`
}

type ValMoniker struct {
	Moniker      string `db:"moniker"`
	OperatorAddr string `db:"operator_addr"`
}

type MaxBlockHeight struct {
	MaxBlockHeightSign   sql.NullInt64 `db:"max_block_height_sign"`
	MaxBlockHeightMissed sql.NullInt64 `db:"max_block_height_missed"`
}

type ValRanking struct {
	ChainName    string
	BlockHeight  int64
	Moniker      string
	OperatorAddr string
	Ranking      int
}
