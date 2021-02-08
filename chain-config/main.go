package main

import (
	"fmt"
	"log"

	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
)

func main() {
	log.Println("simulate app config")
	if !custom.IsSetAppConfig() {
		panic(fmt.Errorf("bech32 is not set corretly"))
	}
}
