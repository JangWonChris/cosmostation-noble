package exporter

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	consensustypes "github.com/tendermint/tendermint/consensus/types"
	tenderminttypes "github.com/tendermint/tendermint/types"
	"go.uber.org/zap"
)

type ConsensusState struct {
	HeightRoundStep   string    `json:"height/round/step"`
	StartTime         time.Time `json:"start_time"`
	ProposalBlockHash string    `json:"proposal_block_hash"`
	LockedBlockHash   string    `json:"locked_block_hash"`
	ValidBlockHash    string    `json:"valid_block_hash"`
	HeightVoteSet     []struct {
		Round              int      `json:"round"`
		Prevotes           []string `json:"prevotes"`
		PrevotesBitArray   string   `json:"prevotes_bit_array"`
		Precommits         []string `json:"precommits"`
		PrecommitsBitArray string   `json:"precommits_bit_array"`
	} `json:"height_vote_set"`
	Proposer struct {
		Address string `json:"address"`
		Index   int    `json:"index"`
	} `json:"proposer"`
}

// RoundStepNewHeight     = RoundStepType(0x01) // Wait til CommitTime + timeoutCommit
// RoundStepNewRound      = RoundStepType(0x02) // Setup new round and go to RoundStepPropose
// RoundStepPropose       = RoundStepType(0x03) // Did propose, gossip proposal
// RoundStepPrevote       = RoundStepType(0x04) // Did prevote, gossip prevotes
// RoundStepPrevoteWait   = RoundStepType(0x05) // Did receive any +2/3 prevotes, start timeout
// RoundStepPrecommit     = RoundStepType(0x06) // Did precommit, gossip precommits
// RoundStepPrecommitWait = RoundStepType(0x07) // Did receive any +2/3 precommits, start timeout
// RoundStepCommit        = RoundStepType(0x08) // Entered commit state machine
// consensus/types/round_state
func (ex *Exporter) getConsensusState() error {

	ctx := context.Background()

	subscriber := "consensus"
	cap := 10
	q := tenderminttypes.QueryForEvent(tenderminttypes.EventNewRound)
	// q := "tm.event = 'NewBlock'"
	out, err := ex.Client.RPC.Subscribe(ctx, subscriber, q.String(), cap)
	if err != nil {
		return err
	}

	go func() {
		for e := range out {
			// zap.S().Info("got ", e.Data.(tenderminttypes.EventDataNewRound))
			zap.S().Info("got ", e.Data)
			zap.S().Info(tenderminttypes.EventQueryNewRound.Matches(e.Events))
		}
	}()
	time.Sleep(30 * time.Second)
	os.Exit(1)
	for {
		select {

		case a := <-out:
			zap.S().Info(a.Data)
			match, err := tenderminttypes.EventQueryPolka.Matches(a.Events)
			if err != nil {
				return err
			}
			zap.S().Info(match)
		}
	}

	rawCS, err := ex.Client.RPC.ConsensusState(ctx)
	if err != nil {
		return err
	}

	cs := &ConsensusState{}

	json.Unmarshal(rawCS.RoundState, cs)

	hrs := strings.Split(cs.HeightRoundStep, "/")
	height, err := strconv.Atoi(hrs[0])
	if err != nil {
		return err
	}
	round, err := strconv.Atoi(hrs[1])
	if err != nil {
		return err
	}
	step, err := strconv.ParseInt(hrs[2], 10, 8)
	if err != nil {
		return err
	}
	_, _, _ = height, round, step

	var msg string
	switch consensustypes.RoundStepType(step) {
	case consensustypes.RoundStepNewHeight: // Wait til CommitTime + timeoutCommit
		msg = "New Height"
	case consensustypes.RoundStepNewRound: // Setup new round and go to RoundStepPropose
		msg = "New Round"
	case consensustypes.RoundStepPropose: // Did propose, gossip proposal
		msg = "Propose"
	case consensustypes.RoundStepPrevote: // Did prevote, gossip prevotes
		msg = "Prevote"
	case consensustypes.RoundStepPrevoteWait: // Did receive any +2/3 prevotes, start timeout
		msg = "Prevote Wait"
	case consensustypes.RoundStepPrecommit: // Did precommit, gossip precommits
		msg = "Precommit"
	case consensustypes.RoundStepPrecommitWait: // Did receive any +2/3 precommits, start timeout
		msg = "Precommit Wait"
	case consensustypes.RoundStepCommit: // Entered commit state machine
		msg = "Commit"
	}
	// zap.S().Info(msg)
	zap.S().Info("height/round/step : ", height, round, step, msg)

	// zap.S().Info(cs.HeightRoundStep)
	// zap.S().Info(cs.HeightVoteSet[round])

	// dumpcs, err := ex.Client.RPC.DumpConsensusState(ctx)
	// if err != nil {
	// 	return err
	// }

	// zap.S().Info(string(rawCS.RoundState))
	// zap.S().Info(dumpcs.RoundState)

	return nil
}
