package exporter

import (
	"log"
)

// Task1
func (ses *StatsExporterService) RunningTask1() {
	log.Println("Task1")
}

// Task2
func (ses *StatsExporterService) RunningTask2() {
	log.Println("Task2------")
}

// Task3
func (ses *StatsExporterService) RunningTask3(name string) {
	log.Println("Task3: ", name)
}
