package models

import (
	"github.com/labstack/gommon/log"
	"stream/database"
	"stream/utils"
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
	BidBTC float64 `json:"bidBTC"`
	AskBTC float64 `json:"askBTC"`
	MidBTC float64 `json:"midBTC"`
	BidSOL float64 `json:"bidSOL"`
	AskSOL float64 `json:"askSOL"`
	MidSOL float64 `json:"midSOL"`
}

var backup []utils.BAM

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
	if len(series) == 2 {
		backup = series
		return Prices{BidBTC: series[0].Bid, AskBTC: series[0].Ask, MidBTC: series[0].Mid, BidSOL: series[1].Bid, AskSOL: series[1].Ask, MidSOL: series[1].Mid}
	} else {
		if len(backup) == 2 {
			return Prices{BidBTC: backup[0].Bid, AskBTC: backup[0].Ask, MidBTC: backup[0].Mid, BidSOL: backup[1].Bid, AskSOL: backup[1].Ask, MidSOL: backup[1].Mid}
		} else {
			return Prices{BidBTC: 0.000, AskBTC: 0.000, MidBTC: 0.000, BidSOL: 0.000, AskSOL: 0.000, MidSOL: 0.000}
		}
	}
}

// WritePrices Generates random values and writes to passed on channel
func WritePrices(ch1 chan<- Prices) error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		newVal := getPrices()
		ch1 <- newVal
	}
}
