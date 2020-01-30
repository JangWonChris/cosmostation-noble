package models

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TxData struct {
	Txs []string `json:"txs"`
}

// type GeneralTx struct {
// 	Height    json.RawMessage `json:"height"`
// 	TxHash    json.RawMessage `json:"txhash"`
//  Data   	  json.RawMessage `json:"data"`
// 	RawLog    json.RawMessage `json:"raw_log"`
// 	Logs      json.RawMessage `json:"logs"`
// 	GasWanted json.RawMessage `json:"gas_wanted"`
// 	GasUsed   json.RawMessage `json:"gas_used"`
// 	Tags      json.RawMessage `json:"tags"`
// 	Tx        json.RawMessage `json:"tx"`
// 	Timestamp json.RawMessage `json:"timestamp"`
// }

// GeneralTx is a struct for general tx
type GeneralTx struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	Data   string `json:"data"`
	RawLog string `json:"raw_log"`
	Logs   []struct {
		MsgIndex string `json:"msg_index"`
		Success  bool   `json:"success"`
		Log      string `json:"log"`
	} `json:"logs"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
	Tags      []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"tags"`
	Tx struct {
		Type  string `json:"type"`
		Value struct {
			Msg []struct {
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			} `json:"msg"`
			Fee struct {
				Amount string `json:"amount"`
				Gas    string `json:"gas"`
			} `json:"fee"`
			Signatures json.RawMessage `json:"signatures"`
			Memo       string          `json:"memo"`
		} `json:"value"`
	} `json:"tx"`
	Timestamp time.Time `json:"timestamp"`
}

type CreateValidatorMsgValueTx struct {
	Description struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	} `json:"description"`
	Commission struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
	} `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
	DelegatorAddress  string `json:"delegator_address"`
	ValidatorAddress  string `json:"validator_address"`
	Pubkey            string `json:"pubkey"`
	Value             struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"value"`
}

type DelegateMsgValueTx struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

type UndelegateMsgValueTx struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

type RedelegateMsgValueTx struct {
	DelegatorAddress    string `json:"delegator_address"`
	ValidatorSrcAddress string `json:"validator_src_address"`
	ValidatorDstAddress string `json:"validator_dst_address"`
	Amount              struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

type SubmitProposalMsgValueTx struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	ProposalType   string `json:"proposal_type"`
	Proposer       string `json:"proposer"`
	InitialDeposit []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"initial_deposit"`
}

type VoteMsgValueTx struct {
	ProposalID string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}

type DepositMsgValueTx struct {
	ProposalID string `json:"proposal_id"`
	Depositor  string `json:"depositor"`
	Amount     []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"amount"`
}
