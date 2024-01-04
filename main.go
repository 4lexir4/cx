package main

import (
	"fmt"
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/4lexir4/cx/server"
)

const ethPrice = 1_281

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	go makeMarketSimple(c)

	select {}
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

		fmt.Println("best bid", bestBid)
		fmt.Println("best ask", bestAsk)

		<-ticker.C
	}
}
