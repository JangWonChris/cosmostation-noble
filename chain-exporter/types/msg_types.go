package types

// Params for validator_set_infos table
const (
	EventTypeMsgCreateValidator = "create_validator"
	EventTypeMsgEditValidator   = "edit_validator"
	EventTypeMsgDelegate        = "delegate"
	EventTypeMsgUndelegate      = "begin_unbonding"
	EventTypeMsgBeginRedelegate = "begin_redelegate"
)
