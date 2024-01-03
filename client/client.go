package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	//"fmt"
	"net/http"

	"github.com/4lexir4/cx/orderbook"
	"github.com/4lexir4/cx/server"
)

const Endpoint = "http://localhost:3000"

type Client struct {
	*http.Client
}

func NewClient() *Client {
	return &Client{
		Client: http.DefaultClient,
	}
}

func (c *Client) GetTrades(market string) ([]*orderbook.Trade, error) {
	e := fmt.Sprintf("%s/trades/%s", Endpoint, market)

	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	trades := []*orderbook.Trade{}

	if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, err
	}

	return trades, nil

}

func (c *Client) GetOrders(userID int64) (*server.GetOrdersResponse, error) {
	e := fmt.Sprintf("%s/order/%d", Endpoint, userID)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	orders := server.GetOrdersResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (c *Client) PlaceMarketOrder(p *PlaceOrderParams) (*server.PlaceOrderResponse, error) {
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.MarketOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Market: server.MarketETH,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	e := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	placeOrderResponse := &server.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, err
	}

	return placeOrderResponse, nil
}

func (c *Client) GetBestBid() (float64, error) {
	e := fmt.Sprintf("%s/book/ETH/bid", Endpoint)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	priceResp := &server.PriceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return 0, err
	}
	return priceResp.Price, err
}

func (c *Client) GetBestAsk() (float64, error) {
	e := fmt.Sprintf("%s/book/ETH/ask", Endpoint)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	priceResp := &server.PriceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return 0, err
	}
	return priceResp.Price, err
}

func (c *Client) CancelOrder(orderID int64) error {
	e := fmt.Sprintf("%s/order/%d", Endpoint, orderID)
	req, err := http.NewRequest(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	if err != nil {
		return err
	}
	return nil
}

type PlaceOrderParams struct {
	UserID int64
	Bid    bool
	// price only needed for placing limit order
	Price float64
	Size  float64
}

func (c *Client) PlaceLimitOrder(p *PlaceOrderParams) (*server.PlaceOrderResponse, error) {
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.LimitOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Price:  p.Price,
		Market: server.MarketETH,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	e := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	placeOrderResponse := &server.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, err
	}

	return placeOrderResponse, nil
}
