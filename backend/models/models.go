package models

import (
	"github.com/labstack/gommon/log"
	"stream/database"
	"time"
)

type currency1 struct {
}

type Currency1 interface {
	// GetLiveCurrency1 Return read only channel that has information about price
	GetLiveCurrency1() (<-chan Prices, error)
}

func NewCurrency1() Currency1 {
	return &currency1{}
}

type Prices struct {
	Bid float64 `json:"bid"` //bid
	Ask float64 `json:"ask"` //mid
	Mid float64 `json:"mid"` //ask
}

// GetLiveCurrency1 Consider this model responsible for puling data from backend and serving it to controller
func (m *currency1) GetLiveCurrency1() (<-chan Prices, error) {
	ch1 := make(chan Prices)
	go func() {
		err := WritePrices(ch1)
		if err != nil {
			log.Error(err)
		}
	}()

	return ch1, nil
}

func getPrices() Prices {
	series := database.Price()
	return Prices{Bid: series.Bid, Ask: series.Ask, Mid: series.Mid}
}

// WritePrices Generates random values and writes to passed on channel
func WritePrices(ch1 chan<- Prices) error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		ch1 <- getPrices()
	}
}
