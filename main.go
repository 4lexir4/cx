package main

import (
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/4lexir4/cx/mm"
	"github.com/4lexir4/cx/server"
)

func main() {
	go server.StartServer()
	time.Sleep(1 * time.Second)

	c := client.NewClient()

	cfg := mm.Config{
		OrderSize:      10,
		MinSpread:      100,
		MakeInterval:   1 * time.Second,
		SeedOffset:     400,
		ExchnageClient: c,
		UserID:         8,
	}
	maker := mm.NewMarketMaker(cfg)

	maker.Start()

	time.Sleep(1 * time.Second)

	go marketOrderPlacer(c)

	select {}
}

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		buyOrder := client.PlaceOrderParams{
			UserID: 7,
			Bid:    true,
			Size:   1,
		}
		_, err := c.PlaceMarketOrder(&buyOrder)
		if err != nil {
			panic(err)
		}
		sellOrder := client.PlaceOrderParams{
			UserID: 7,
			Bid:    false,
			Size:   1,
		}
		_, err = c.PlaceMarketOrder(&sellOrder)
		if err != nil {
			panic(err)
		}

		<-ticker.C
	}
}
