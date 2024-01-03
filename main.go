package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/4lexir4/cx/server"
)

const (
	userID = 7
)

var (
	tick = 2 * time.Second
)

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(5 * time.Second)
	for {

		trades, err := c.GetTrades("ETH")
		if err != nil {
			panic(err)
		}

		if len(trades) > 0 {
			fmt.Printf("EXCHANGE PRICE: %.2f\n", trades[len(trades)-1].Price)
		}

		otherMarketSellOrder := &client.PlaceOrderParams{
			UserID: 8,
			Bid:    false,
			Size:   1000,
		}

		orderResp, err := c.PlaceMarketOrder(otherMarketSellOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketSellOrder := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    false,
			Size:   100,
		}

		orderResp, err = c.PlaceMarketOrder(marketSellOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketBuyOrder := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    true,
			Size:   100,
		}

		orderResp, err = c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		<-ticker.C
	}
}

func makeMarketSimple(c *client.Client) {
	ticker := time.NewTicker(tick)

	for {
		orders, err := c.GetOrders(userID)
		if err != nil {
			log.Println(err)
		}

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

		if len(orders.Bids) < 3 {
			// place the bid order
			bidLimit := &client.PlaceOrderParams{
				UserID: 7,
				Bid:    true,
				Price:  bestBid + 100,
				Size:   1_000,
			}

			bidOrderResp, err := c.PlaceLimitOrder(bidLimit)
			if err != nil {
				log.Println(bidOrderResp.OrderID)
			}
		}

		if len(orders.Asks) < 3 {
			// place the ask order
			askLimit := &client.PlaceOrderParams{
				UserID: 7,
				Bid:    false,
				Price:  bestAsk - 100,
				Size:   1_000,
			}

			askOrderResp, err := c.PlaceLimitOrder(askLimit)
			if err != nil {
				log.Println(askOrderResp.OrderID)
			}
		}

		fmt.Println("Best ask price", bestAsk)
		fmt.Println("Best bid price", bestBid)

		<-ticker.C
	}

}

func seedMarket(c *client.Client) error {
	ask := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  10_000,
		Size:   10_000,
	}
	_, err := c.PlaceLimitOrder(ask)
	if err != nil {
		return err
	}

	bid := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  9_000,
		Size:   10_000,
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

	go makeMarketSimple(c)
	time.Sleep(1 * time.Second)
	marketOrderPlacer(c)

	select {}
}
