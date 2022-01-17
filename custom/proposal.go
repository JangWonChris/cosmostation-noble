package custom

import (
	"fmt"
	"log"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func GetProposalType(c govtypes.Content) (proposalType string) {
	switch i := c.(type) {
	default:
		log.Printf("unrecognized proposal type : %T\n", i)
		proposalType = fmt.Sprintf("%T", i)
	}
	return
}
