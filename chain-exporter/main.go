package main

import (
	"flag"
	"log"
	"os"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/exporter"
	"go.uber.org/zap"
)

func main() {
	/*
		chain-exporter
			-mode=basic (default)
			-mode=raw
			-mode=refine
			-mode=genesis
	*/
	mode := flag.String("mode", "basic", "chain-exporter mode \n basic : default, will store current chain status\n raw : will only store jsonRawMessage of block and transaction to database\n refine : refine new data from database the legacy chain stored\n genesis : extract genesis state from the given file")
	initialHeight := flag.Int64("initial-height", 0, "initial height of chain-exporter to sync")
	gPath := flag.String("genesis-file-path", "", "absolute path of genesis.json")

	//deprecated => is replaced with mode=genesis
	genesisParse := flag.Bool("genesis-parse", false, "if true, chain-exporter will parse genesis state and store database")
	//deprecated
	chunkOnly := flag.Bool("chunk-only", false, "if true, chain-exporter will only store chunks of blocks and transactions")

	flag.Parse()

	log.Println("mode : ", *mode)
	log.Println("path : ", *gPath)
	log.Println("parse : ", *genesisParse)
	log.Println("initial-height :", *initialHeight)
	log.Println("chunk-only :", *chunkOnly)

	// parse인데, path가 같이 안나오면 오류
	// 2안 parse 옵션이 있고, path가 없으면 default path를 알아서 잡고, 오류 시 프로그램 종료
	// if *genesis-parse && *gPath == "" {
	// 	log.Fatal("genesis-file-path must be defined")
	// 	os.Exit(1)
	// }

	os.Exit(0)

	exporter := exporter.NewExporter()

	if *genesisParse {
		if err := exporter.GetGenesisStateFromGenesisFile(*gPath); err != nil {
			os.Exit(1)
		}
		zap.S().Info("parse complete")
	}

	exporter.Start(*initialHeight, *chunkOnly)
}
