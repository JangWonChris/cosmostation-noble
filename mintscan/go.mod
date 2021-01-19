module github.com/cosmostation/cosmostation-cosmos/mintscan

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.40.0
	github.com/cosmos/gaia v1.0.0
	github.com/cosmos/gaia/v3 v3.0.0
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210119170926-440efb4265df
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/go-resty/resty/v2 v2.3.0
	github.com/gorilla/mux v1.8.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.1
	github.com/tendermint/tm-db v0.6.3
	github.com/test-go/testify v1.1.4 // indirect
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.33.2
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
