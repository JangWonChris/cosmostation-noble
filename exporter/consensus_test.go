package exporter

import (
	"testing"
	"time"
)

func TestGetConsensusState(t *testing.T) {

	// ex.getConsensusState()
	// os.Exit(1)
	ex.Client.RPC.Start()
	cnt := 0

	ex.getConsensusState()
	for {
		time.Sleep(300 * time.Millisecond)
		if cnt > 100 {
			ex.Client.RPC.Stop()
			break
		}
		cnt++
	}

}
