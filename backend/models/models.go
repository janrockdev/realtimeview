package models

import (
	"fmt"
	lr "github.com/sirupsen/logrus"
	kdb "github.com/sv/kdbgo"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"runtime"
	"strconv"
	"stream/utils"
	"strings"
	"time"
)

var logr = &lr.Logger{
	Out:   os.Stdout,
	Level: lr.DebugLevel,
	Formatter: &prefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	},
}

var res []utils.BAM

type currency1 struct {
}

type Currency1 interface {
	// GetLiveCurrency1 Return read only channel that has information about price
	GetLiveCurrency1() (<-chan utils.Prices, error)
}

func NewCurrency1() Currency1 {
	return &currency1{}
}

// GetLiveCurrency1 Consider this model responsible for puling data from backend and serving it to controller
func (m *currency1) GetLiveCurrency1() (<-chan utils.Prices, error) {
	ch1 := make(chan utils.Prices)
	go func() {
		logr.Debug("New session started.")
		WriteData(ch1)
	}()

	return ch1, nil
}

// toStruct wrap KDB output to buy-ask-mid struct
func toStruct(tbl kdb.Table) []utils.BAM {

	var data []utils.BAM

	nrows := tbl.Data[0].Len()
	for i := 0; i < nrows; i++ {
		rec := utils.BAM{Time: tbl.Data[0].Index(i).(time.Time), Bid: tbl.Data[1].Index(i).(float64), Ask: tbl.Data[2].Index(i).(float64), Mid: tbl.Data[3].Index(i).(float64)}
		data = append(data, rec)
	}
	return data
}

// routineId return goroutine-thread ID
func routineId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func WriteData(ch1 chan<- utils.Prices) {

	con, _ := kdb.DialKDB("192.168.0.89", 5000, "")

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool, 1)

	go func() {

		var counter int

		for {
			select {
			case <-done:
				logr.Infof("Timeout for goroutine: %v", routineId())
				return
			case _ = <-ticker.C:

				counter += 1

				logr.Infof("Gofunc: %v | loop: %v", routineId(), counter)

				res = []utils.BAM{}
				ktbl, err := con.Call("select [-1] ts, bid, ask, mid, mid, mid from quotes where sym=`BTCUSD")
				if err != nil {
					fmt.Println("Query failed:", err)
					return
				}

				series := toStruct(ktbl.Data.(kdb.Table))

				//for _, v := range series {
				//	res = v.Open
				//}

				res = append(res, series[0])

				ktblSOL, err := con.Call("select [-1] ts, bid, ask, mid, mid, mid from quotes where sym=`SOLUSD")
				seriesSOL := toStruct(ktblSOL.Data.(kdb.Table))

				res = append(res, seriesSOL[0])

				logr.Infof("Result(%v): %v", counter, res)

				var backup []utils.BAM

				if len(res) == 2 {
					backup = series
					ch1 <- utils.Prices{BidBTC: res[0].Bid, AskBTC: res[0].Ask, MidBTC: res[0].Mid, BidSOL: res[1].Bid, AskSOL: res[1].Ask, MidSOL: res[1].Mid}
				} else {
					if len(backup) == 2 {
						ch1 <- utils.Prices{BidBTC: backup[0].Bid, AskBTC: backup[0].Ask, MidBTC: backup[0].Mid, BidSOL: backup[1].Bid, AskSOL: backup[1].Ask, MidSOL: backup[1].Mid}
					} else {
						ch1 <- utils.Prices{BidBTC: 0.000, AskBTC: 0.000, MidBTC: 0.000, BidSOL: 0.000, AskSOL: 0.000, MidSOL: 0.000}
					}
				}

				// Exit strategy for goroutine
				if counter == 10 { //5sec * 12 * 15 = 900
					logr.Info("Terminating...")
					done <- true
				}
			}
		}
	}()

	return
}
