package exporter

import (
	"os"
	"testing"
)

var (
	exporter *Exporter
)

func TestMain(m *testing.M) {
	exporter = NewExporter()

	os.Exit(m.Run())
}

func TestExportAllStats(t *testing.T) {
	exporter.SaveStatsMarket5M()
	exporter.SaveStatsMarket1H()
	exporter.SaveStatsMarket1D()

	exporter.SaveNetworkStats1H()
	exporter.SaveNetworkStats1D()

	exporter.SaveValidatorsStats1H()
	exporter.SaveValidatorsStats1D()
}

func TestExportMarketStats(t *testing.T) {
	exporter.SaveStatsMarket5M()
	exporter.SaveStatsMarket1H()
	exporter.SaveStatsMarket1D()
}

func TestExportNetworkStats(t *testing.T) {
	exporter.SaveNetworkStats1H()
	exporter.SaveNetworkStats1D()
}

func TestExportValidatorStats(t *testing.T) {
	exporter.SaveValidatorsStats1H()
	exporter.SaveValidatorsStats1D()
}
