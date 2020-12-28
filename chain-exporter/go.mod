module github.com/cosmostation/cosmostation-cosmos/chain-exporter

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.40.0-rc5.0.20201223042253-92317931c21f
	github.com/cosmos/gaia/v3 v3.0.0-20201223164621-fe60afe07059
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/go-resty/resty/v2 v2.3.0
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.0
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.33.2
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
