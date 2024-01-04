package mm

import (
	"time"

	"github.com/4lexir4/cx/client"
	"github.com/sirupsen/logrus"
)

type Config struct {
	UserID         int64
	OrderSize      float64
	MinSpread      float64
	SeedOffset     float64
	ExchnageClient *client.Client
	MakeInterval   time.Duration
	PriceOffset    float64
}

type MarketMaker struct {
	userID         int64
	orderSize      float64
	minSpread      float64
	seedOffset     float64
	exchnageClient *client.Client
	makeInterval   time.Duration
	priceOffset    float64
}

func NewMarketMaker(cfg Config) *MarketMaker {
	return &MarketMaker{
		userID:         cfg.UserID,
		orderSize:      cfg.OrderSize,
		minSpread:      cfg.MinSpread,
		seedOffset:     cfg.SeedOffset,
		exchnageClient: cfg.ExchnageClient,
		makeInterval:   cfg.MakeInterval,
		priceOffset:    cfg.PriceOffset,
	}
}

func (mm *MarketMaker) Start() {
	logrus.WithFields(logrus.Fields{
		"id":           mm.userID,
		"orderSize":    mm.orderSize,
		"makeInterval": mm.makeInterval,
		"minSpread":    mm.minSpread,
		"priceOffset":  mm.priceOffset,
	}).Info("starting market maker")
	go mm.makerLoop()
}

func (mm *MarketMaker) makerLoop() {
	ticker := time.NewTicker(mm.makeInterval)

	for {
		bestBid, err := mm.exchnageClient.GetBestBid()
		if err != nil {
			logrus.Error(err)
			break
		}

		bestAsk, err := mm.exchnageClient.GetBestAsk()
		if err != nil {
			logrus.Error(err)
			break
		}

		if bestAsk.Price == 0 && bestBid.Price == 0 {
			if err := mm.SeedMarket(); err != nil {
				logrus.Error(err)
				break
			}
			continue
		}

		if bestBid.Price == 0 {
			bestBid.Price = bestAsk.Price - mm.priceOffset*2
		}

		if bestAsk.Price == 0 {
			bestAsk.Price = bestBid.Price + mm.priceOffset*2
		}

		spread := bestAsk.Price - bestBid.Price

		if spread <= mm.minSpread {
			continue
		}

		if err := mm.placeOrder(true, bestBid.Price+mm.priceOffset); err != nil {
			logrus.Error(err)
			break
		}

		if err := mm.placeOrder(false, bestAsk.Price-mm.priceOffset); err != nil {
			logrus.Error(err)
			break
		}
		<-ticker.C
	}
}

func (mm *MarketMaker) placeOrder(bid bool, price float64) error {
	bidOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Bid:    bid,
		Price:  price,
		Size:   mm.orderSize,
	}
	_, err := mm.exchnageClient.PlaceLimitOrder(bidOrder)
	return err
}

func (mm *MarketMaker) SeedMarket() error {
	currentPrice := simulateFetchCurrentETHPrice()

	logrus.WithFields(logrus.Fields{
		"currentETHPrice": currentPrice,
		"seedOffset":      mm.seedOffset,
	}).Info("orderbooks empty >>> seeding market")
	priceOffset := mm.seedOffset

	bidOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Bid:    true,
		Price:  currentPrice - priceOffset,
		Size:   mm.orderSize,
	}
	_, err := mm.exchnageClient.PlaceLimitOrder(bidOrder)
	if err != nil {
		return err
	}
	askOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Bid:    false,
		Price:  currentPrice + priceOffset,
		Size:   mm.minSpread,
	}
	_, err = mm.exchnageClient.PlaceLimitOrder(askOrder)
	return err
}

// this simulates a call to another exchange fetching
// the current ETH price so we can offset both for bid and ask
func simulateFetchCurrentETHPrice() float64 {
	// simulate http round trip
	time.Sleep(80 * time.Millisecond)

	return 1000.0
}
