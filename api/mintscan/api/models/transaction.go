package models

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ResultTransactionInfo struct {
	Height int64     `json:"height"`
	TxHash string    `json:"tx_hash"`
	Time   time.Time `json:"time"`
}

type GeneralTx struct {
	Height    json.RawMessage `json:"height"`
	TxHash    json.RawMessage `json:"txhash"`
	RawLog    json.RawMessage `json:"raw_log"`
	Code      json.RawMessage `json:"code"`
	Logs      json.RawMessage `json:"logs"`
	GasWanted json.RawMessage `json:"gas_wanted"`
	GasUsed   json.RawMessage `json:"gas_used"`
	Tags      json.RawMessage `json:"tags"`
	Tx        json.RawMessage `json:"tx"`
	Timestamp json.RawMessage `json:"timestamp"`
}

/*
	LCD - Transaction for Power Event - Delegate, Redelegate, Unbonding, Create Validator
	Tx API - Type 별 파싱하기 위해 준비했던 흔적들인데 위에 json.RawMessage로 받아주니 쉽게 해결 됬다.
	혹시 모르니 지금은 남겨둔다
*/

type (
	// https://lcd.cosmostation.io/txs/284184911ba7d486f752c2354576bebb3d4780d98608be49830a48ac573f06af
	DelegateMsgValueTx struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Value            struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"value"`
	}
	// https://lcd.cosmostation.io/txs/1ABCE9FE606B5CBC991E09B1D2EA0263EAE34FF49B19E09E0330F6CA165B2C96
	UndelegateMsgValueTx struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Amount           struct {
			Denom  string  `json:"denom"`
			Amount sdk.Dec `json:"amount"`
		} `json:"amount"`
	}
	// https://lcd.cosmostation.io/txs/48853f0ff5d23096b90f3844ea97938936eae651625f4b3109ca5e7406d91edf
	RedelegateMsgValueTx struct {
		DelegatorAddress    string `json:"delegator_address"`
		ValidatorSrcAddress string `json:"validator_src_address"`
		ValidatorDstAddress string `json:"validator_dst_address"`
		Amount              struct {
			Denom  string  `json:"denom"`
			Amount sdk.Dec `json:"amount"`
		} `json:"amount"`
	}
	// https://lcd.cosmostation.io/txs/038acdb508fcec44f19c6f8048176dc329dac6da4a2fcc6e37a1499fd0d9b491
	CreateValidatorMsgValueTx struct {
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
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"value"`
	}
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

type DelegateLCD struct {
	Height string `json:"height"`
	Txhash string `json:"txhash"`
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
				Type  string `json:"type"`
				Value struct {
					DelegatorAddress string `json:"delegator_address"`
					ValidatorAddress string `json:"validator_address"`
					Value            struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"value"`
				} `json:"value"`
			} `json:"msg"`
			Fee struct {
				Amount []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"amount"`
				Gas string `json:"gas"`
			} `json:"fee"`
			Signatures []struct {
				PubKey struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"pub_key"`
				Signature string `json:"signature"`
			} `json:"signatures"`
			Memo string `json:"memo"`
		} `json:"value"`
	} `json:"tx"`
}
