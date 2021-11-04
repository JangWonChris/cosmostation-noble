module github.com/cosmostation/cosmostation-cosmos

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.42.9
	github.com/cosmos/gaia/v5 v5.0.5
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210630022738-4a764a0f0cee
	github.com/cosmostation/mintscan-database v0.0.0-20211104111540-af9a15513ba2
	github.com/cosmostation/mintscan-prometheus v0.0.0-20210628093844-2404f3c78830
	github.com/go-pg/pg/v10 v10.9.3
	github.com/go-resty/resty/v2 v2.4.0
	github.com/gorilla/mux v1.8.0
	github.com/gravity-devs/liquidity v1.2.9
	github.com/prometheus/client_golang v1.11.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
