package types

const (
	// DefaultQueryValidatorsPage is the default page number for querying validators via querier.
	DefaultQueryValidatorsPage = 1

	// DefaultQueryValidatorsPerPage is the default per page number for querying validators via querier.
	DefaultQueryValidatorsPerPage = 200

	// BondedValidatorStatus is status code when a validator is live.
	BondedValidatorStatus = 2

	// UnbondingValidatorStatus is status code when a validator is not live.
	UnbondingValidatorStatus = 1

	// UnbondedValidatorStatus is status code when a validator is jailed.
	UnbondedValidatorStatus = 0
)
