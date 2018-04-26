package main

import (
	"fmt"
	"log"

	cmc "github.com/miguelmota/go-coinmarketcap"
)

func main() {
	// Get global market data
	marketInfo, err := cmc.GetGlobalMarketData()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(marketInfo)
}
