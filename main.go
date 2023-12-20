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

	bidParams := &client.PlaceLimitOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  10_000,
		Size:   1000,
	}
	go func() {
		for {
			resp, err := c.PlaceLimitOrder(bidParams)
			if err != nil {
				panic(err)
			}

			fmt.Println("Order ID =>", resp.OrderID)

			if err := c.CancelOrder(resp.OrderID); err != nil {
				panic(err)
			}

			time.Sleep(1 * time.Second)
		}
	}()
	askParams := &client.PlaceLimitOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  8_000,
		Size:   1000,
	}

	for {
		resp, err := c.PlaceLimitOrder(askParams)
		if err != nil {
			panic(err)
		}

		fmt.Println("Order ID =>", resp.OrderID)

		time.Sleep(1 * time.Second)
	}
	select {}
}
