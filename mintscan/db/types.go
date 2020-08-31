package db

// Params that are used when querying transactions that are made by an account
const (
	// Params that are in Cosmos SDK
	QueryTxParamFromAddress      = "messages->0->'value'->>'from_address' = "
	QueryTxParamToAddress        = "messages->0->'value'->>'to_address' = "
	QueryTxParamInputsAddress    = "messages->0->'value'->'inputs'->0->>'address' = "
	QueryTxParamOutpusAddress    = "messages->0->'value'->'outputs'->0->>'address' = "
	QueryTxParamDelegatorAddress = "messages->0->'value'->>'delegator_address' = "
	QueryTxParamAddress          = "messages->0->'value'->>'address' = "
	QueryTxParamProposer         = "messages->0->'value'->>'proposer' = "
	QueryTxParamDepositer        = "messages->0->'value'->>'depositor' = "
	QueryTxParamVoter            = "messages->0->'value'->>'voter' = "

	// Params that are in Cosmos SDK and they related to validators
	QueryTxParamValidatorAddress    = "messages->0->'value'->>'validator_address' = "
	QueryTxParamValidatorDstAddress = "messages->0->'value'->>'validator_dst_address' = "
	QueryTxParamValidatorSrcAddress = "messages->0->'value'->>'validator_src_address' = "
	QueryTxParamValidatorCommission = "messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission'"

	// Params that are related to DeFi
	QueryTxParamFrom      = "messages->0->'value'->>'from' = "
	QueryTxParamSender    = "messages->0->'value'->>'sender' = "
	QueryTxParamOwner     = "messages->0->'value'->>'owner' = "
	QueryTxParamDepositor = "messages->0->'value'->>'depositor' = "

	// Param that is for specific denom
	QueryTxParamDenom = "messages->0->'value'->'amount'->0->>'denom' = "
)
