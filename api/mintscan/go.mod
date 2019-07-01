module github.com/cosmostation/cosmostation-cosmos/api/mintscan

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.34.3
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/gorilla/mux v1.7.3
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/spf13/viper v1.4.0
	github.com/tendermint/tendermint v0.31.5
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	gopkg.in/resty.v1 v1.12.0
	mellium.im/sasl v0.2.1 // indirect
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
