package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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

type PlaceLimitOrderParams struct {
	UserID int64
	Bid    bool
	Price  float64
	Size   float64
}

func (c *Client) PlaceLimitOrder(p *PlaceLimitOrderParams) error {
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
		return err
	}

	e := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", resp)

	return nil
}
