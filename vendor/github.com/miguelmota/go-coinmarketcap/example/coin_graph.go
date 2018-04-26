package main

import (
	"fmt"
	"log"
	"time"

	cmc "github.com/miguelmota/go-coinmarketcap"
)

func main() {
	threeMonths := int64(60 * 60 * 24 * 90)
	now := time.Now()
	secs := now.Unix()
	start := secs - threeMonths
	end := secs

	// Get graph data for coin
	coinGraphData, err := cmc.GetCoinGraphData("ethereum", start, end)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(coinGraphData)
}
