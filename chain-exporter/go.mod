module github.com/cosmostation/cosmostation-cosmos/chain-exporter

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.37.4
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/go-resty/resty/v2 v2.3.0
	github.com/golang/mock v1.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.32.7
	go.uber.org/zap v1.13.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	mellium.im/sasl v0.2.1 // indirect
)

// replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
