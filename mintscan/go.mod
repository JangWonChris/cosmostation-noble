module github.com/cosmostation/cosmostation-cosmos/mintscan

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.41.3
	github.com/cosmos/gaia/v4 v4.0.4
	github.com/cosmostation/cosmostation-cosmos/chain-config v0.0.0-00010101000000-000000000000
	// github.com/cosmostation/mintscan-backend-library v0.0.0-20210215124422-0da1f2875834
	// github.com/cosmostation/mintscan-backend-library v0.0.0-20210221065353-c439d341db6d
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210222091607-09fabc04bacb
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/go-resty/resty/v2 v2.4.0
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.7.0
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.35.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmostation/cosmostation-cosmos/chain-config => ../chain-config

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
