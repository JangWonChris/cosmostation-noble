package types

// parameters for validator_set_infos table
const (
	EventTypeMsgCreateValidator = "create_validator"
	EventTypeMsgEditValidator   = "edit_validator"
	EventTypeMsgDelegate        = "delegate"
	EventTypeMsgUndelegate      = "begin_unbonding"
	EventTypeMsgBeginRedelegate = "begin_redelegate"
)

type Signature struct {
	Address   string `json:"address,omitempty"`
	Pubkey    string `json:"pubkey,omitempty"`
	Signature string `json:"signature,omitempty"`
}
