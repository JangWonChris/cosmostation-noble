package db

// Params that are used when querying transactions that are made by an account.
const (
	// Params that are used for general transactions in Cosmos SDK.
	QueryTxParamFromAddress      = "messages->0->>'from_address' = "
	QueryTxParamToAddress        = "messages->0->>'to_address' = "
	QueryTxParamInputsAddress    = "messages->0->'inputs'->0->>'address' = "
	QueryTxParamOutpusAddress    = "messages->0->'outputs'->0->>'address' = "
	QueryTxParamDelegatorAddress = "messages->0->>'delegator_address' = "
	QueryTxParamAddress          = "messages->0->>'address' = "
	QueryTxParamProposer         = "messages->0->>'proposer' = "
	QueryTxParamDepositer        = "messages->0->>'depositor' = "
	QueryTxParamVoter            = "messages->0->>'voter' = "

	// Params that are used for validators in Cosmos SDK.
	QueryTxParamValidatorAddress    = "messages->0->>'validator_address' = "
	QueryTxParamValidatorDstAddress = "messages->0->>'validator_dst_address' = "
	QueryTxParamValidatorSrcAddress = "messages->0->>'validator_src_address' = "
	QueryTxParamValidatorCommission = "messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission'"

	// Param that are used for selecting particular coin denomination.
	// QueryTxParamDenom = "messages->0->'value'->'amount'->0->>'denom' = "
	QueryTxParamDenom = "messages->0->'amount'->0->>'denom' = "
)
