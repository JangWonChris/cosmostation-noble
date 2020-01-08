package db

import (
	"log"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"
)

// InsertCoinGeckoMarket1H saves StatsCoingeckoMarket1H
func (db *Database) InsertCoinGeckoMarket1H(data schema.StatsCoingeckoMarket1H) (bool, error) {
	_, err := db.Model(&data).Insert()
	if err != nil {
		log.Printf("failed to insert CoinGeckoMarkek1H: %v \n", err)
		return false, nil
	}
	return true, nil
}

// InsertCoinGeckoMarket24H saves StatsCoingeckoMarket24H
func (db *Database) InsertCoinGeckoMarket24H(data schema.StatsCoingeckoMarket24H) (bool, error) {
	_, err := db.Model(&data).Insert()
	if err != nil {
		log.Printf("failed to insert CoinGeckoMarkek24H: %v \n", err)
		return false, nil
	}
	return true, nil
}

// InsertCoinMarketCapMarket1H saves StatsCoinmarketcapMarket1H
func (db *Database) InsertCoinMarketCapMarket1H(data schema.StatsCoinmarketcapMarket1H) (bool, error) {
	_, err := db.Model(&data).Insert()
	if err != nil {
		log.Printf("failed to insert CoinmarketcapMarket1H: %v \n", err)
		return false, nil
	}
	return true, nil
}

// InsertCoinMarketCapMarket24H saves StatsCoinmarketcapMarket24H
func (db *Database) InsertCoinMarketCapMarket24H(data schema.StatsCoinmarketcapMarket24H) (bool, error) {
	_, err := db.Model(&data).Insert()
	if err != nil {
		log.Printf("failed to insert CoinmarketcapMarket1H: %v \n", err)
		return false, nil
	}
	return true, nil
}

// InsertNetworkStats1H saves StatsNetwork1H
func (db *Database) InsertNetworkStats1H(data schema.StatsNetwork1H) (bool, error) {
	_, err := db.Model(&data).Insert()
	if err != nil {
		log.Printf("failed to insert StatsNetwork1H: %v \n", err)
		return false, nil
	}
	return true, nil
}

// InsertNetworkStats24H saves StatsNetwork24H
func (db *Database) InsertNetworkStats24H(data schema.StatsNetwork24H) (bool, error) {
	_, err := db.Model(&data).Insert()
	if err != nil {
		log.Printf("failed to insert StatsNetwork24H: %v \n", err)
		return false, nil
	}
	return true, nil
}
