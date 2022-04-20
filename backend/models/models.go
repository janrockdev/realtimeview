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
	GetLiveCurrency1() (<-chan MidPrice, error)
}

func NewCurrency1() Currency1 {
	return &currency1{}
}

type MidPrice struct {
	Price float64 `json:"midPrice"`
}

// GetLiveCurrency1 Consider this model responsible for puling data from backend and serving it to controller
func (m *currency1) GetLiveCurrency1() (<-chan MidPrice, error) {
	ch1 := make(chan MidPrice)
	go func() {
		err := WritePrices(ch1)
		if err != nil {
			log.Error(err)
		}
	}()

	return ch1, nil
}

// WritePrices Generates random values and writes to passed on channel
func WritePrices(ch1 chan<- MidPrice) error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		newPrice := MidPrice{Price: database.Price()}
		ch1 <- newPrice
	}
}
