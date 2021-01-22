module github.com/cosmostation/cosmostation-cosmos/mintscan

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.40.0
	github.com/cosmos/gaia/v3 v3.0.0
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210120031210-8342355c7571
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/go-resty/resty/v2 v2.3.0
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.6.1
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.33.2
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
