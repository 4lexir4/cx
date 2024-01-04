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
}

type MarketMaker struct {
	userID         int64
	orderSize      float64
	minSpread      float64
	seedOffset     float64
	exchnageClient *client.Client
	makeInterval   time.Duration
}

func NewMarketMaker(cfg Config) *MarketMaker {
	return &MarketMaker{
		userID:         cfg.UserID,
		orderSize:      cfg.OrderSize,
		minSpread:      cfg.MinSpread,
		seedOffset:     cfg.SeedOffset,
		exchnageClient: cfg.ExchnageClient,
		makeInterval:   cfg.MakeInterval,
	}
}

func (mm *MarketMaker) Start() {
	logrus.WithFields(logrus.Fields{
		"id":           mm.userID,
		"orderSize":    mm.orderSize,
		"makeInterval": mm.makeInterval,
		"minSpread":    mm.minSpread,
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

		if bestAsk == 0 && bestBid == 0 {
			if err := mm.SeedMarket(); err != nil {
				logrus.Error(err)
				break
			}
		}
		<-ticker.C
	}
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
