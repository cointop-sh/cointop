package main

import (
	"fmt"
	"log"

	cmc "github.com/miguelmota/go-coinmarketcap"
)

func main() {
	// Get info about coin
	coinInfo, err := cmc.GetCoinData("ethereum")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(coinInfo)
}
