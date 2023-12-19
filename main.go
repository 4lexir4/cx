package main

import (
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/4lexir4/cx/server"
)

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	params := &client.PlaceLimitOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  4000.0,
		Size:   4000.0,
	}

	if err := c.PlaceLimitOrder(params); err != nil {
		panic(err)
	}
}
