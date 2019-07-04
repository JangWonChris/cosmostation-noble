module github.com/cosmostation/cosmostation-cosmos/chain-exporter

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.34.7
	github.com/cosmostation/cosmostation-kava/chain-exporter v0.0.0-20190701133858-6108846bcb3b
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/olivere/elastic v6.2.21+incompatible
	github.com/spf13/viper v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.5
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
	gopkg.in/resty.v1 v1.12.0
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
