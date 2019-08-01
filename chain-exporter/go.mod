module github.com/cosmostation/cosmostation-cosmos/chain-exporter

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.36.0-rc1
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/spf13/viper v1.4.0
	github.com/tendermint/tendermint v0.32.1
	gopkg.in/resty.v1 v1.12.0
	mellium.im/sasl v0.2.1 // indirect
)

// replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
