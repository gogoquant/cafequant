package main

import (
	"fmt"
	"os"
	"snack.com/xiyanxiyan10/conver"
	dbtypes "snack.com/xiyanxiyan10/stockdb/types"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

// main ...
func main() {
	if len(os.Args) < 3 {
		fmt.Println("命令行的参数不合法:", len(os.Args))
		return
	}
	if err := config.Init(os.Args[1]); err != nil {
		fmt.Printf("config init error is %s\n", err.Error())
		return
	}
	var opt constant.Option
	var constract = "quarter"
	var symbol = "BTC/USD"
	var io = "online"
	var period = "H1"
	var path = os.Args[2]

	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.LogBack = true

	maker := api.ExchangeMaker[opt.Type]
	exchange, err := maker(opt)
	if err != nil {
		fmt.Printf("init exchange fail:%s\n", err.Error())
		return
	}

	exchange.SetIO(io)
	exchange.SetContractType(constract)
	exchange.SetStockType(symbol)
	exchange.SetPeriod(period)
	exchange.SetPeriodSize(10)
	exchange.Start()

	vec, err := util.ReadCSV(path)
	if err != nil {
		fmt.Printf("read csv fail:%s\n", err.Error())
		return
	}
	//fmt.Printf("%s\n", util.Struct2Json(vec))
	var datums []dbtypes.OHLC
	for i, data := range vec {
		var ohlc dbtypes.OHLC
		if i == 0 {
			continue
		}

		ohlc.Time = conver.Int64Must(data[0])
		ohlc.Open = conver.Float64Must(data[1])
		ohlc.High = conver.Float64Must(data[2])
		ohlc.Low = conver.Float64Must(data[3])
		ohlc.Close = conver.Float64Must(data[4])
		ohlc.Volume = conver.Float64Must(data[5]) * ohlc.Close
		datums = append(datums, ohlc)
	}
	err = exchange.BackPutOHLCs(datums, exchange.GetPeriod())

	if err != nil {
		fmt.Printf("put ohlc fail:%s\n", err.Error())
		return
	}
}
