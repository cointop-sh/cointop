package coinmarketcap

import (
	"testing"
	"time"
)

func TestGetGlobalMarketData(t *testing.T) {
	market, err := GetGlobalMarketData()
	if err != nil {
		t.FailNow()
	}

	if market.ActiveAssets == 0 {
		t.FailNow()
	}
	if market.ActiveCurrencies == 0 {
		t.FailNow()
	}
	if market.ActiveMarkets == 0 {
		t.FailNow()
	}
	if market.BitcoinPercentageOfMarketCap == 0 {
		t.FailNow()
	}
	if market.Total24HVolumeUSD == 0 {
		t.FailNow()
	}
	if market.TotalMarketCapUSD == 0 {
		t.FailNow()
	}
}

func TestGetGlobalMarketGraphData(t *testing.T) {
	var threeMonths int64 = (60 * 60 * 24 * 90)
	end := time.Now().Unix()
	start := end - threeMonths

	graph, err := GetGlobalMarketGraphData(start, end)
	if err != nil {
		t.FailNow()
	}

	if graph.MarketCapByAvailableSupply[0][0] == 0 {
		t.FailNow()
	}

	if graph.VolumeUSD[0][0] == 0 {
		t.FailNow()
	}
}

func TestGetAltcoinMarketGraphData(t *testing.T) {
	var threeMonths int64 = (60 * 60 * 24 * 90)
	end := time.Now().Unix()
	start := end - threeMonths

	graph, err := GetAltcoinMarketGraphData(start, end)
	if err != nil {
		t.FailNow()
	}

	if graph.MarketCapByAvailableSupply[0][0] == 0 {
		t.FailNow()
	}

	if graph.VolumeUSD[0][0] == 0 {
		t.FailNow()
	}
}

func TestGetCoinData(t *testing.T) {
	coin, err := GetCoinData("ethereum")
	if err != nil {
		t.FailNow()
	}

	if coin.AvailableSupply == 0 {
		t.FailNow()
	}
	if coin.ID == "" {
		t.FailNow()
	}
	if coin.LastUpdated == "" {
		t.FailNow()
	}
	if coin.MarketCapUSD == 0 {
		t.FailNow()
	}
	if coin.Name == "" {
		t.FailNow()
	}
	if coin.PercentChange1H == 0 {
		t.FailNow()
	}
	if coin.PercentChange24H == 0 {
		t.FailNow()
	}
	if coin.PercentChange7D == 0 {
		t.FailNow()
	}
	if coin.PriceBTC == 0 {
		t.FailNow()
	}
	if coin.PriceUSD == 0 {
		t.FailNow()
	}
	if coin.Rank == 0 {
		t.FailNow()
	}
	if coin.Symbol == "" {
		t.FailNow()
	}
	if coin.TotalSupply == 0 {
		t.FailNow()
	}
	if coin.USD24HVolume == 0 {
		t.FailNow()
	}
}

func TestGetAllCoinData(t *testing.T) {
	coins, err := GetAllCoinData(10)
	if err != nil {
		t.FailNow()
	}

	if len(coins) != 10 {
		t.FailNow()
	}
}

func TestGetCoinGraphData(t *testing.T) {
	var threeMonths int64 = (60 * 60 * 24 * 90)
	end := time.Now().Unix()
	start := end - threeMonths

	graph, err := GetCoinGraphData("ethereum", start, end)
	if err != nil {
		t.FailNow()
	}

	if graph.MarketCapByAvailableSupply[0][0] == 0 {
		t.FailNow()
	}
	if graph.PriceBTC[0][0] == 0 {
		t.FailNow()
	}
	if graph.PriceUSD[0][0] == 0 {
		t.FailNow()
	}
	if graph.VolumeUSD[0][0] == 0 {
		t.FailNow()
	}
}

func TestGetCoinPriceUSD(t *testing.T) {
	price, err := GetCoinPriceUSD("ethereum")
	if err != nil {
		t.FailNow()
	}
	if price <= 0 {
		t.FailNow()
	}
}

func TestGetCoinMarkets(t *testing.T) {
	markets, err := GetCoinMarkets("ethereum")
	if err != nil {
		t.FailNow()
	}
	if len(markets) == 0 {
		t.FailNow()
	}

	market := markets[0]
	if market.Rank == 0 {
		t.FailNow()
	}
	if market.Exchange == "" {
		t.FailNow()
	}
	if market.Pair == "" {
		t.FailNow()
	}
	if market.VolumeUSD == 0 {
		t.FailNow()
	}
	if market.Price == 0 {
		t.FailNow()
	}
	if market.VolumePercent == 0 {
		t.FailNow()
	}
	if market.Updated == "" {
	}
}

func TestDoReq(t *testing.T) {
	// TODO
}

func TestMakeReq(t *testing.T) {
	// TODO
}

func TestToInt(t *testing.T) {
	v := toInt("5")
	if v != 5 {
		t.FailNow()
	}
}

func TestToFloat(t *testing.T) {
	v := toFloat("5.2")
	if v != 5.2 {
		t.FailNow()
	}
}
