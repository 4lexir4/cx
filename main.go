package main

import (
	"math/rand"
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
		MinSpread:      20,
		MakeInterval:   1 * time.Second,
		SeedOffset:     40,
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
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		randInt := rand.Intn(10)
		bid := true
		if randInt > 5 {
			bid = false
		}

		order := client.PlaceOrderParams{
			UserID: 7,
			Bid:    bid,
			Size:   1,
		}
		_, err := c.PlaceMarketOrder(&order)
		if err != nil {
			panic(err)
		}

		<-ticker.C
	}
}
