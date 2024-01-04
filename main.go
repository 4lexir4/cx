package main

import (
	"fmt"
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/4lexir4/cx/server"
)

const ethPrice = 1_281.0

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	go makeMarketSimple(c)

	select {}
}

func seedMarket(c *client.Client) {
	// NOTE this should be an async call to fetch the price
	currentPrice := ethPrice

	priceOffset := 100.0

	bidOrder := client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  currentPrice - priceOffset,
		Size:   10,
	}
	_, err := c.PlaceLimitOrder(&bidOrder)
	if err != nil {
		panic(err)
	}
	askOrder := client.PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  currentPrice + priceOffset,
		Size:   10,
	}
	_, err = c.PlaceLimitOrder(&askOrder)
	if err != nil {
		panic(err)
	}
}

func makeMarketSimple(c *client.Client) {
	ticker := time.NewTicker(1 * time.Second)
	for {

		bestAsk, err := c.GetBestAsk()
		if err != nil {
			panic(err)
		}

		bestBid, err := c.GetBestBid()
		if err != nil {
			panic(err)
		}

		if bestAsk == 0 && bestBid == 0 {
			seedMarket(c)
			continue
		}

		fmt.Println("best bid", bestBid)
		fmt.Println("best ask", bestAsk)

		<-ticker.C
	}
}
