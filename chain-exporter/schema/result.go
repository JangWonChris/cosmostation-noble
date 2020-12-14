package schema

// ExportData has all exported data that will be saved in database.
type ExportData struct {
	ResultAccounts []Account
	ResultBlock    Block
	// ResultGenesisAccounts             []Account
	ResultTxs                         []TransactionLegacy
	ResultTxsMessages                 []TransactionAccount
	ResultEvidence                    []Evidence
	ResultMissBlocks                  []Miss
	ResultMissDetailBlocks            []MissDetail
	ResultAccumulatedMissBlocks       []Miss
	ResultProposals                   []Proposal
	ResultDeposits                    []Deposit
	ResultVotes                       []Vote
	ResultGenesisValidatorsSet        []PowerEventHistory
	ResultValidatorsPowerEventHistory []PowerEventHistory
}

type ExportRawData struct {
	ResultBlockJSONChunk RawBlock
	ResultTxsJSONChunk   []RawTransaction
}

// NewExportData returns a new ExportData.
func NewExportData(e ExportData) *ExportData {
	return &ExportData{
		ResultAccounts: e.ResultAccounts,
		ResultBlock:    e.ResultBlock,
		// ResultGenesisAccounts:             e.ResultGenesisAccounts,
		ResultTxs:                         e.ResultTxs,
		ResultEvidence:                    e.ResultEvidence,
		ResultMissBlocks:                  e.ResultMissBlocks,
		ResultMissDetailBlocks:            e.ResultMissDetailBlocks,
		ResultAccumulatedMissBlocks:       e.ResultAccumulatedMissBlocks,
		ResultProposals:                   e.ResultProposals,
		ResultDeposits:                    e.ResultDeposits,
		ResultVotes:                       e.ResultVotes,
		ResultGenesisValidatorsSet:        e.ResultGenesisValidatorsSet,
		ResultValidatorsPowerEventHistory: e.ResultValidatorsPowerEventHistory,
	}
}
