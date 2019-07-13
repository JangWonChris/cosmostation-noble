package main

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
)

func main() {
	fmt.Println("Stats Exporter")

	// Configuration in config.yaml
	config := config.NewConfig()
}
