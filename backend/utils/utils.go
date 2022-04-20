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
