package coinmarketcap

import "testing"

func TestV1Tickers(t *testing.T) {
	coins, err := V1Tickers(10, "EUR")
	if err != nil {
		t.FailNow()
	}

	if len(coins) != 10 {
		t.FailNow()
	}
}
