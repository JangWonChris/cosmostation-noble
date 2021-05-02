module github.com/cosmostation/cosmostation-cosmos/chain-config

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/gaia/v4 v4.2.1
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210208045014-5ba1778df744
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.13.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
