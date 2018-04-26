package main

import (
	"fmt"
	"log"

	cmc "github.com/miguelmota/go-coinmarketcap"
)

func main() {
	// get data for all coins
	coins, err := cmc.GetAllCoinData(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, coin := range coins {
		fmt.Println(coin.Symbol, coin.PriceUSD)
	}
}
