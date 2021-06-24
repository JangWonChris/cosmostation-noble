module github.com/cosmostation/cosmostation-cosmos

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/gaia/v4 v4.2.1
	// github.com/cosmostation/mintscan-backend-library v0.0.0-20210609122107-607af66ead19
	// github.com/cosmostation/mintscan-backend-library v0.0.0-20210615083719-9ce78896626e
	github.com/cosmostation/mintscan-backend-library v0.0.0-20210616095652-cd86c0e1cd86
	github.com/cosmostation/mintscan-database v0.0.0-20210624053100-4110b2442284 //v0.0.0-20210615051646-f80c2a1fdf9e
	github.com/cosmostation/mintscan-prometheus v0.0.0-20210616095618-b14862ad83a2
	github.com/go-pg/pg/v10 v10.9.1
	github.com/go-resty/resty/v2 v2.4.0
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.9
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.35.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
