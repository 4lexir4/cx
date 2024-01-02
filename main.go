package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/4lexir4/cx/server"
)

var tick = 2 * time.Second

func makeMarketSimple(c *client.Client) {
	ticker := time.NewTicker(tick)

	for {
		<-ticker.C

		bestAsk, err := c.GetBestAsk()
		if err != nil {
			log.Println(err)
		}
		bestBid, err := c.GetBestBid()
		if err != nil {
			log.Println(err)
		}

		spread := math.Abs(bestBid - bestAsk)
		fmt.Println("Spread", spread)

		bidLimit := &client.PlaceOrderParams{
			UserID: 7,
			Bid:    true,
			Price:  bestBid + 100,
			Size:   1_000,
		}

		bidOrderID, err := c.PlaceLimitOrder(bidLimit)
		if err != nil {
			log.Println(bidOrderID)
		}

		askLimit := &client.PlaceOrderParams{
			UserID: 7,
			Bid:    false,
			Price:  bestAsk - 100,
			Size:   1_000,
		}

		askOrderID, err := c.PlaceLimitOrder(askLimit)
		if err != nil {
			log.Println(askOrderID)
		}

		fmt.Println("Best ask price", bestAsk)
		fmt.Println("Best bid price", bestBid)
	}

}

func seedMarket(c *client.Client) error {
	ask := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  10_000,
		Size:   1_000_000,
	}
	_, err := c.PlaceLimitOrder(ask)
	if err != nil {
		return err
	}

	bid := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  9_000,
		Size:   1_000_000,
	}
	_, err = c.PlaceLimitOrder(bid)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	if err := seedMarket(c); err != nil {
		panic(err)
	}

	makeMarketSimple(c)

	//for {
	//	limitOrderParams := &client.PlaceOrderParams{
	//		UserID: 8,
	//		Bid:    false,
	//		Price:  10_000,
	//		Size:   500_000,
	//	}
	//	_, err := c.PlaceLimitOrder(limitOrderParams)
	//	if err != nil {
	//		panic(err)
	//	}

	//	otherLimitOrderParams := &client.PlaceOrderParams{
	//		UserID: 666,
	//		Bid:    false,
	//		Price:  9_000,
	//		Size:   5_000_000,
	//	}
	//	_, err = c.PlaceLimitOrder(otherLimitOrderParams)
	//	if err != nil {
	//		panic(err)
	//	}

	//	buyLimitOrderParams := &client.PlaceOrderParams{
	//		UserID: 666,
	//		Bid:    true,
	//		Price:  11_000,
	//		Size:   5_000_000,
	//	}
	//	_, err = c.PlaceLimitOrder(buyLimitOrderParams)
	//	if err != nil {
	//		panic(err)
	//	}
	//	//fmt.Println("Placed limit order from the client =>", resp.OrderID)

	//	marketOrderParams := &client.PlaceOrderParams{
	//		UserID: 7,
	//		Bid:    true,
	//		Size:   1_000_000,
	//	}
	//	_, err = c.PlaceMarketOrder(marketOrderParams)
	//	if err != nil {
	//		panic(err)
	//	}

	//	bestBidPrice, err := c.GetBestBid()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("Best bid price: [%.2f]\n", bestBidPrice)

	//	bestAskPrice, err := c.GetBestAsk()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("Best ask price: [%.2f]\n", bestAskPrice)

	//	time.Sleep(1 * time.Second)
	//}
	select {}
}
