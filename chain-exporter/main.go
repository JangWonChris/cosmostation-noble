package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/custom"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter"
	"go.uber.org/zap"
)

func init() {
	custom.SetAppConfig()
	if !custom.IsSetBech32() {
		panic(fmt.Errorf("bech32 is not set corretly"))
	}
	log.Println("Current bech32 : ", sdktypes.GetConfig())
}

func main() {
	mode := flag.String("mode", "basic", "chain-exporter mode \n  - basic : default, will store current chain status\n  - raw : will only store jsonRawMessage of block and transaction to database\n  - refine : refine new data from database the legacy chain stored\n  - genesis : extract genesis state from the given file")
	initialHeight := flag.Int64("initial-height", 0, "initial height of chain-exporter to sync")
	genesisFilePath := flag.String("genesis-file-path", "", "absolute path of genesis.json")
	flag.Parse()

	log.Println("mode : ", *mode)
	log.Println("genesis-file-path : ", *genesisFilePath)
	log.Println("initial-height :", *initialHeight)

	op := exporter.BASIC_MODE

	switch *mode {
	case "basic": //기본 동작
		op = exporter.BASIC_MODE
	case "raw":
		op = exporter.RAW_MODE
	case "refine":
		op = exporter.REFINE_MODE
	case "genesis":
		op = exporter.GENESIS_MODE
	default:
		log.Println("Unknow operator type :", *mode)
		os.Exit(1)
	}

	ex := exporter.NewExporter()

	if op == exporter.GENESIS_MODE {
		if err := ex.GetGenesisStateFromGenesisFile(*genesisFilePath); err != nil {
			os.Exit(1)
		}
		zap.S().Info("genesis file parsing complete")
		os.Exit(0)
	}

	if op == exporter.REFINE_MODE {
		if err := ex.Refine(op); err != nil {
			os.Exit(1)
		}
		zap.S().Info("refine successfully complete")
		os.Exit(0)
	}

	ex.Start(*initialHeight, op)
}
