package exporter

import (
	"log"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
)

// Task1
func (ses *ChainExporterService) RunningTask1() {
	log.Println("Task1")
}

// Task2
func (ses *ChainExporterService) RunningTask2(config *config.Config) {
	log.Println("Task2------", config.db.Host)
}

// Task3
func (ses *ChainExporterService) RunningTask3(name string) {
	log.Println("Task3: ", name)
}
