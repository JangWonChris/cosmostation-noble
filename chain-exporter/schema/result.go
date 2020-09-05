package schema

// ExportData has all exported data that will be saved in database.
type ExportData struct {
	ResultBlock                       Block
	ResultTxs                         []Transaction
	ResultEvidence                    []Evidence
	ResultMissBlocks                  []Miss
	ResultMissDetailBlocks            []MissDetail
	ResultAccumulatedMissBlocks       []Miss
	ResultProposals                   []Proposal
	ResultDeposits                    []Deposit
	ReusltVotes                       []Vote
	ResultGenesisValidatorsSet        []PowerEventHistory
	ResultValidatorsPowerEventHistory []PowerEventHistory
}

// NewExportData returns a new ExportData.
func NewExportData(e ExportData) *ExportData {
	return &ExportData{
		ResultBlock:                       e.ResultBlock,
		ResultTxs:                         e.ResultTxs,
		ResultEvidence:                    e.ResultEvidence,
		ResultMissBlocks:                  e.ResultMissBlocks,
		ResultMissDetailBlocks:            e.ResultMissDetailBlocks,
		ResultAccumulatedMissBlocks:       e.ResultAccumulatedMissBlocks,
		ResultProposals:                   e.ResultProposals,
		ResultDeposits:                    e.ResultDeposits,
		ReusltVotes:                       e.ReusltVotes,
		ResultGenesisValidatorsSet:        e.ResultGenesisValidatorsSet,
		ResultValidatorsPowerEventHistory: e.ResultValidatorsPowerEventHistory,
	}
}
