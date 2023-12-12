package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/4lexir4/cx/orderbook"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
)

const (
	MarketOrder        OrderType = "MARKET"
	LimitOrder         OrderType = "LIMIT"
	exchangePrivateKey           = "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
	MarketETH          Market    = "ETH"
)

type (
	OrderType string

	Market string

	Exchange struct {
		PrivateKey *ecdsa.PrivateKey
		orderbooks map[Market]*orderbook.Orderbook
	}

	PlaceOrderRequest struct {
		Type   OrderType // limit or market
		Bid    bool
		Size   float64
		Price  float64
		Market Market
	}

	Order struct {
		ID        int64
		Price     float64
		Size      float64
		Bid       bool
		Timestamp int64
	}

	OrderbookData struct {
		TotalBidVolume float64
		TotalAskVolume float64
		Asks           []*Order
		Bids           []*Order
	}

	MatchedOrder struct {
		Price float64
		Size  float64
		ID    int64
	}
)

func main() {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	ex, err := NewExchange(exchangePrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceOrder)
	e.DELETE("/order/:id", ex.cancelOrder)

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA("4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)

	gasLimit := uint64(21000) // in units

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e")

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID := big.NewInt(1337)
	//chainID, err := client.NetworkID(context.Background()) // 1337
	//if err != nil {
	//	log.Fatal(err)
	//}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%+v", tx)

	balance, err := client.BalanceAt(ctx, toAddress, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance)

	e.Start(":3000")
}

type User struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewUser(privateKey string) *User {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	return &User{
		PrivateKey: pk,
	}
}

func httpErrorHandler(err error, c echo.Context) {
	fmt.Println(err)
}

func NewExchange(privateKey string) (*Exchange, error) {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		PrivateKey: pk,
		orderbooks: orderbooks,
	}, nil
}

func (ex *Exchange) handleGetBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found"})
	}

	orderbookData := OrderbookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*Order{},
		Bids:           []*Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Bids = append(orderbookData.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}

func (ex *Exchange) cancelOrder(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	ob := ex.orderbooks[MarketETH]
	order := ob.Orders[int64(id)]
	ob.CancelOrder(order)

	return c.JSON(200, map[string]any{"msg": "order deleted"})
}

func (ex *Exchange) handlePlaceMarketOrder(market Market, order *orderbook.Order) ([]orderbook.Match, []*MatchedOrder) {
	ob := ex.orderbooks[market]
	matches := ob.PlaceMarketOrder(order)
	matchedOrders := make([]*MatchedOrder, len(matches))

	isBid := false
	if order.Bid {
		isBid = true
	}

	for i := 0; i < len(matchedOrders); i++ {
		id := matches[i].Bid.ID
		if isBid {
			id = matches[i].Ask.ID
		}
		matchedOrders[i] = &MatchedOrder{
			ID:    id,
			Size:  matches[i].SizeFilled,
			Price: matches[i].Price,
		}
	}
	return matches, matchedOrders
}

func (ex *Exchange) handlePlaceLimitOrder(market Market, price float64, order *orderbook.Order) error {
	ob := ex.orderbooks[market]
	ob.PlaceLimitOrder(price, order)
	return nil
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)

	if placeOrderData.Type == LimitOrder {
		if err := ex.handlePlaceLimitOrder(market, placeOrderData.Price, order); err != nil {
			return err
		}
		return c.JSON(200, map[string]any{"msg": "limit order placed"})
	}

	if placeOrderData.Type == MarketOrder {
		matches, matchedOrders := ex.handlePlaceMarketOrder(market, order)
		if err := ex.handleMatches(matches); err != nil {
			return err
		}
		return c.JSON(200, map[string]any{"matches": matchedOrders})
	}

	return nil
}

func (ex *Exchange) handleMatches(matches []orderbook.Match) error {
	return nil
}
