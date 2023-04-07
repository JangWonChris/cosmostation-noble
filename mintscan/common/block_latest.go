package common

import (
	"net/http"

	"github.com/cosmostation/cosmostation-noble/app"
	"github.com/cosmostation/cosmostation-noble/errors"
	"go.uber.org/zap"
)

type NodeInfo struct {
	Network           string `json:"network"`
	LatestBlockHeight int64  `json:"latest_block_height"`
}

func GetBlocksLatest(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		status, err := a.Client.RPC.GetStatus()
		if err != nil {
			zap.S().Debug("failed to get network status", zap.Error(err))
			errors.ErrServerUnavailable(rw, http.StatusInternalServerError)
			return
		}
		nodeInfo := NodeInfo{
			Network:           status.NodeInfo.Network,
			LatestBlockHeight: status.SyncInfo.LatestBlockHeight,
		}

		respond(rw, nodeInfo)
	}
}
