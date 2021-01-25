module github.com/cosmostation/cosmostation-cosmos/chain-exporter

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.40.0
	github.com/cosmos/gaia/v3 v3.0.0
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210125133527-87e11a16ee6a
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/go-resty/resty/v2 v2.3.0
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.1
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.33.2
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
