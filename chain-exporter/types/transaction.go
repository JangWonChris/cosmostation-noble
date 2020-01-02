package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GeneralTx is general tx struct that is unmarshallable for any tx_msg type
type GeneralTx struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	RawLog string `json:"raw_log"`
	Logs   []struct {
		MsgIndex uint8  `json:"msg_index"`
		Success  bool   `json:"success"`
		Log      string `json:"log"`
	} `json:"logs"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
	Events    []struct {
		Type       string `json:"type"`
		Attributes []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
	} `json:"events"`
	Tx struct {
		Type  string `json:"type"`
		Value struct {
			Msg []struct {
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			} `json:"msg"`
			Fee struct {
				Amount []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"amount"`
				Gas string `json:"gas"`
			} `json:"fee"`
			Signatures json.RawMessage `json:"signatures"`
			Memo       string          `json:"memo"`
		} `json:"value"`
	} `json:"tx"`
	Timestamp time.Time `json:"timestamp"`
}

type MsgSend struct {
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	Amount      sdk.Coins `json:"amount"`
}

type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs" yaml:"inputs"`
	Outputs []Output `json:"outputs" yaml:"outputs"`
}

// Input models transaction input
type Input struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
}

// Output models transaction outputs
type Output struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
}

type MsgCreateValidator struct {
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

type MsgDelegate struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

type MsgUndelegate struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

type MsgBeginRedelegate struct {
	DelegatorAddress    string `json:"delegator_address"`
	ValidatorSrcAddress string `json:"validator_src_address"`
	ValidatorDstAddress string `json:"validator_dst_address"`
	Amount              struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

type MsgSubmitProposal struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	ProposalType   string `json:"proposal_type"`
	Proposer       string `json:"proposer"`
	InitialDeposit []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"initial_deposit"`
}

type MsgVote struct {
	ProposalID string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}

type MsgDeposit struct {
	ProposalID string `json:"proposal_id"`
	Depositor  string `json:"depositor"`
	Amount     []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"amount"`
}
