package utils

import (
	"time"
)

//BAM is a data structure
type BAM struct {
	Time time.Time
	Bid  float64
	Ask  float64
	Mid  float64
}

type Prices struct {
	BidBTC float64 `json:"bidBTC"`
	AskBTC float64 `json:"askBTC"`
	MidBTC float64 `json:"midBTC"`
	BidSOL float64 `json:"bidSOL"`
	AskSOL float64 `json:"askSOL"`
	MidSOL float64 `json:"midSOL"`
}
