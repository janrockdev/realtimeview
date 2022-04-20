package models

import (
	"math/rand"
	"time"
)

type model struct {
}

type Model interface {
	// Return read only channel that has information about % current cpu
	GetLiveCpuUsage() (<-chan Cpu, error)
}

func NewModel() Model {
	return &model{}
}

type Cpu struct {
	PercentageUsage int `json:"cpuPercentageUsage"`
}

// Consider this model responsible for puling data from backend and serving it to controller
func (m *model) GetLiveCpuUsage() (<-chan Cpu, error) {
	ch := make(chan Cpu)
	go WriteValues(ch)

	return ch, nil
}

// Generates random values and writes to passed on channel
func WriteValues(ch chan<- Cpu) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		rand.Seed(time.Now().UnixNano())
		min := 1980
		max := 2000
		newVal := Cpu{PercentageUsage: rand.Intn(max-min+1) + min}
		ch <- newVal
	}
}

//// Generates random values and writes to passed on channel
//func WriteValues(ch chan <- Cpu ) error {
//	ticker := time.NewTicker(1 * time.Second)
//	for{
//		<- ticker.C
//		newVal := Cpu{PercentageUsage:rand.Intn(100)}
//		ch <- newVal
//	}
//}
