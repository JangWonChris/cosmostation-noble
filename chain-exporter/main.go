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

	log.Println(*gPath)
	log.Println(*parse)
	log.Println(*initialHeight)

	// parse인데, path가 같이 안나오면 오류
	// 2안 parse 옵션이 있고, path가 없으면 default path를 알아서 잡고, 오류 시 프로그램 종료
	// if *parse && *gPath == "" {
	// 	log.Fatal("genesis-file-path must be defined")
	// 	os.Exit(1)
	// }

	exporter := exporter.NewExporter()

	if *parse {
		if err := exporter.GetGenesisStateFromGenesisFile(*gPath); err != nil {
			os.Exit(1)
		}
		zap.S().Info("parse complete")
	}

	exporter.Start(*initialHeight)
}
