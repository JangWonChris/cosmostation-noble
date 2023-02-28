package exporter

import (
	"testing"
	"time"

	"github.com/cosmostation/cosmostation-noble/app"
)

func TestProposalAlarm(t *testing.T) {
	chainEx := app.NewApp("chain-exporter")
	ex = NewExporter(chainEx)
	go ex.ProposalNotificationToSlack(51)
	time.Sleep(5 * time.Second)
}
