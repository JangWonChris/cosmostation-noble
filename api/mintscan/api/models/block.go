package models

import "time"

type ResultBlock struct {
	Height          int64     `json:"height"`
	Proposer        string    `json:"proposer"`
	OperatorAddress string    `json:"operator_address"`
	Moniker         string    `json:"moniker"`
	BlockHash       string    `json:"block_hash"`
	Identity        string    `json:"identity"`
	NumTxs          int64     `json:"num_txs"`
	TxData          TxData    `json:"tx_data"`
	Time            time.Time `json:"time"`
}

type TxData struct {
	Txs []string `json:"txs"`
}

type ResultBlocksByOperatorAddr struct {
	Height                 int64     `json:"height"`
	Proposer               string    `json:"proposer"`
	OperatorAddress        string    `json:"operator_address"`
	Moniker                string    `json:"moniker"`
	BlockHash              string    `json:"block_hash"`
	Identity               string    `json:"identity"`
	NumTxs                 int64     `json:"num_txs"`
	TotalNumProposerBlocks int       `json:"total_num_proposer_blocks"`
	TxData                 TxData    `json:"tx_data"`
	Time                   time.Time `json:"time"`
}

type ResultBlockByHeight struct {
	Height            int64                       `json:"height"`
	Time              time.Time                   `json:"time"`
	NumTxs            int64                       `json:"num_txs"`
	Proposer          ResultBlockByHeightProposer `json:"proposer"`
	MissingValidators []struct {
		Address         []string `json:"address"`
		OperatorAddress string   `json:"operator_address"`
		Moniker         string   `json:"moniker"`
		VotingPower     float64  `json:"voting_power"`
	} `json:"missing_validators"`
	Evidence []struct {
		Address string `json:"address"`
	} `json:"evidence"`
	TxData TxData `json:"tx_data"`
}

type ResultBlockByHeightProposer struct {
	Address         string  `json:"address"`
	OperatorAddress string  `json:"operator_address"`
	Moniker         string  `json:"moniker"`
	VotingPower     float64 `json:"voting_power"`
}
