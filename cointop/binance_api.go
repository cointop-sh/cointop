package main

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
)

func main() {

	// Auth Part
	var (
		apiKey    = "sJpK9ooVd60yhhFmV2Er7FXnI3AKLpE38UJ20qXvTxnc74yRw19kuqC0yPhGNq1q"
		secretKey = "kujQnkruJeFgI64YKqCLPg3R3sucRqv83V4yUisNoqqnQtUO4ouu7MEKf5lvscwx"
	)
	binance.UseTestnet = true
	client := binance.NewClient(apiKey, secretKey)

	// Using NewListPrices from go-binance lib
	prices, err := client.NewListPricesService().Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range prices {
		fmt.Printf(p.Symbol)
		fmt.Printf("\n")
		fmt.Println(p.Price)
	}
}
