package app

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmostation-cosmos/chain-exporter-es/elastic"
)

type Context struct {
	//Logger        logrus.FieldLogger
	ElasticSearch *elastic.ElasticSearch
	Logger log.Logger
	//RemoteAddress string
}

func (ctx *Context) WithLogger(logger log.Logger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}
