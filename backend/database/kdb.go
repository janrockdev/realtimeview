package database

import (
	"fmt"
	"os"

	"stream/utils"

	"time"

	kdb "github.com/sv/kdbgo"

	lr "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type result struct {
	hlc float64
	ema []float64
	rsi []float64
	bb  []float64
}

var res float64

var logr = &lr.Logger{
	Out:   os.Stdout,
	Level: lr.InfoLevel,
	Formatter: &prefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	},
}

func toStruct(tbl kdb.Table) []utils.OHLCV {

	var data = []utils.OHLCV{}

	nrows := int(tbl.Data[0].Len())
	for i := 0; i < nrows; i++ {
		rec := utils.OHLCV{Time: tbl.Data[0].Index(i).(time.Time), Open: tbl.Data[1].Index(i).(float64), High: tbl.Data[2].Index(i).(float64), Low: tbl.Data[3].Index(i).(float64), Close: tbl.Data[4].Index(i).(float64), Volume: tbl.Data[5].Index(i).(float64)}
		data = append(data, rec)
	}
	return data
}

func Price() float64 {

	con, _ := kdb.DialKDB("192.168.0.89", 5000, "")

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {

		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:

				ktbl, err := con.Call("select [-1] ts, bid, ask, mid, mid, mid from quotes")
				if err != nil {
					fmt.Println("Query failed:", err)
					return
				}

				series := toStruct(ktbl.Data.(kdb.Table))

				for _, v := range series {
					res = v.Open
				}

			}
		}
	}()
	return res
}
