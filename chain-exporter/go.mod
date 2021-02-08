module github.com/cosmostation/cosmostation-cosmos/chain-exporter

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.41.0
	github.com/cosmos/gaia/v4 v4.0.0
	github.com/cosmostation/cosmostation-cosmos/chain-config v0.0.0-00010101000000-000000000000
	// github.com/cosmostation/mintscan-backend-library v0.0.0-20210125133527-87e11a16ee6a
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210208045014-5ba1778df744
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/go-resty/resty/v2 v2.3.0
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.3
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.35.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmostation/cosmostation-cosmos/chain-config => ../chain-config
