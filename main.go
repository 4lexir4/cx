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
			Size:   1_000_000,
		}
		resp, err := c.PlaceLimitOrder(limitOrderParams)
		if err != nil {
			panic(err)
		}
		fmt.Println("Placed limit order from the client =>", resp.OrderID)

		marketOrderParams := &client.PlaceOrderParams{
			UserID: 7,
			Bid:    true,
			Size:   1_000_000,
		}
		resp, err = c.PlaceMarketOrder(marketOrderParams)
		if err != nil {
			panic(err)
		}
		fmt.Println("Placed market order from the client =>", resp.OrderID)

		time.Sleep(1 * time.Second)
	}
	select {}
}
