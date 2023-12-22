package main

import (
	"fmt"
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/4lexir4/cx/server"
)

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	for {
		limitOrderParams := &client.PlaceOrderParams{
			UserID: 8,
			Bid:    false,
			Price:  10_000,
			Size:   500_000,
		}
		_, err := c.PlaceLimitOrder(limitOrderParams)
		if err != nil {
			panic(err)
		}

		otherLimitOrderParams := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    false,
			Price:  9_000,
			Size:   5_000_000,
		}
		_, err = c.PlaceLimitOrder(otherLimitOrderParams)
		if err != nil {
			panic(err)
		}

		buyLimitOrderParams := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    true,
			Price:  11_000,
			Size:   5_000_000,
		}
		_, err = c.PlaceLimitOrder(buyLimitOrderParams)
		if err != nil {
			panic(err)
		}
		//fmt.Println("Placed limit order from the client =>", resp.OrderID)

		marketOrderParams := &client.PlaceOrderParams{
			UserID: 7,
			Bid:    true,
			Size:   1_000_000,
		}
		_, err = c.PlaceMarketOrder(marketOrderParams)
		if err != nil {
			panic(err)
		}

		bestBidPrice, err := c.GetBestBid()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Best bid price: [%.2f]\n", bestBidPrice)

		bestAskPrice, err := c.GetBestAsk()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Best ask price: [%.2f]\n", bestAskPrice)

		time.Sleep(1 * time.Second)
	}
	select {}
}
