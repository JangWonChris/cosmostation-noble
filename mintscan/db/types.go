package db

// Params that are used when querying transactions that are made by an account.
const (
	// Params that are used for general transactions in Cosmos SDK.
	QueryTxParamFromAddress      = "messages->0->'value'->>'from_address' = "
	QueryTxParamToAddress        = "messages->0->'value'->>'to_address' = "
	QueryTxParamInputsAddress    = "messages->0->'value'->'inputs'->0->>'address' = "
	QueryTxParamOutpusAddress    = "messages->0->'value'->'outputs'->0->>'address' = "
	QueryTxParamDelegatorAddress = "messages->0->'value'->>'delegator_address' = "
	QueryTxParamAddress          = "messages->0->'value'->>'address' = "
	QueryTxParamProposer         = "messages->0->'value'->>'proposer' = "
	QueryTxParamDepositer        = "messages->0->'value'->>'depositor' = "
	QueryTxParamVoter            = "messages->0->'value'->>'voter' = "

	// Params that are used for validators in Cosmos SDK.
	QueryTxParamValidatorAddress    = "messages->0->'value'->>'validator_address' = "
	QueryTxParamValidatorDstAddress = "messages->0->'value'->>'validator_dst_address' = "
	QueryTxParamValidatorSrcAddress = "messages->0->'value'->>'validator_src_address' = "
	QueryTxParamValidatorCommission = "messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission'"

	// Param that are used for selecting particular coin denomination.
	QueryTxParamDenom = "messages->0->'value'->'amount'->0->>'denom' = "
)
