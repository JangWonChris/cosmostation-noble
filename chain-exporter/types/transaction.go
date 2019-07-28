package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
	LCD - Transaction for Power Event - Delegate, Redelegate, Unbonding, Create Validator
*/

// GeneralTx is general tx struct that is unmarshallable for any tx_msg type
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

// cosmoshub-1 https://lcd.cosmostation.io/txs/D867403BDD544A0FE41EDB6AE0368CACAADFF6A98BF3AD3E76D868D104DEC2F6
// cosmoshub-2 https://lcd.cosmostation.io/txs/
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

// cosmoshub-1 https://lcd.cosmostation.io/txs/284184911ba7d486f752c2354576bebb3d4780d98608be49830a48ac573f06af
// cosmoshub-2 https://lcd.cosmostation.io/txs/F341DD756A93D5BA8433E92A7D92922BB99A138C85EFCDA8256E958A6E388589
type MsgDelegate struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

// cosmoshub-1 https://lcd.cosmostation.io/txs/1ABCE9FE606B5CBC991E09B1D2EA0263EAE34FF49B19E09E0330F6CA165B2C96
// cosmoshub-2 https://lcd.cosmostation.io/txs/31CD1C12F1DCECDB8C5FE00A01F68E696D2338BCDD9D08AA7702F56F6ACF19FC
type MsgUndelegate struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

// cosmoshub-1 https://lcd.cosmostation.io/txs/48853f0ff5d23096b90f3844ea97938936eae651625f4b3109ca5e7406d91edf
// cosmoshub-2 https://lcd.cosmostation.io/txs/44D490262824A6BD9E110E82586E208E587CB6079E51FBCF450C7B49F2FECA94
type MsgBeginRedelegate struct {
	DelegatorAddress    string `json:"delegator_address"`
	ValidatorSrcAddress string `json:"validator_src_address"`
	ValidatorDstAddress string `json:"validator_dst_address"`
	Amount              struct {
		Denom  string  `json:"denom"`
		Amount sdk.Dec `json:"amount"`
	} `json:"amount"`
}

type (
	// https://lcd.cosmostation.io/txs/A126E48228271FBAEDF49028E4CB724049E09C831769AE402BBB80CCA197C62D
	SubmitProposalMsgValueTx struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		ProposalType   string `json:"proposal_type"`
		Proposer       string `json:"proposer"`
		InitialDeposit []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"initial_deposit"`
	}

	// https://lcd.cosmostation.io/txs/8D17DC38DE754B544F1183AC96FD91D7E9559893A12FCD013E1A87A619856C61
	VoteMsgValueTx struct {
		ProposalID string `json:"proposal_id"`
		Voter      string `json:"voter"`
		Option     string `json:"option"`
	}

	// https://lcd.cosmostation.io/txs/5EED165DE065D07B6772DD0994F8CD23177C5FDE76235012865C5110E52FAF31
	DepositMsgValueTx struct {
		ProposalID string `json:"proposal_id"`
		Depositor  string `json:"depositor"`
		Amount     []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"amount"`
	}
)
