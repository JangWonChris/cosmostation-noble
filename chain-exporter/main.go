package main

import (
	"flag"
	"log"
	"os"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter"
	"go.uber.org/zap"
)

func main() {
	initialHeight := flag.Int64("initial-height", 0, "initial height of chain-exporter to sync")
	gPath := flag.String("genesis-file-path", "", "absolute path of genesis.json")
	parse := flag.Bool("parse", false, "if true, chain-exporter will parse genesis state and store database")

	flag.Parse()

	if *gPath == "" {
		zap.S().Fatal("genesis-file-path must be defined")
		os.Exit(1)
	}

	log.Println(*gPath)
	log.Println(*parse)
	log.Println(*initialHeight)
	exporter := exporter.NewExporter()

	if *parse {
		exporter.GetGenesisStateFromGenesisFile(*gPath)
		zap.S().Info("parse complete")
	}

	exporter.Start(*initialHeight)
}
