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

var res []utils.BAM

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

func toStruct(tbl kdb.Table) []utils.BAM {

	var data = []utils.BAM{}

	nrows := int(tbl.Data[0].Len())
	for i := 0; i < nrows; i++ {
		rec := utils.BAM{Time: tbl.Data[0].Index(i).(time.Time), Bid: tbl.Data[1].Index(i).(float64), Ask: tbl.Data[2].Index(i).(float64), Mid: tbl.Data[3].Index(i).(float64)}
		data = append(data, rec)
	}
	return data
}

func Price() []utils.BAM {

	con, _ := kdb.DialKDB("192.168.0.89", 5000, "")

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {

		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:

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
				return
			}
		}
	}()
	return res
}
