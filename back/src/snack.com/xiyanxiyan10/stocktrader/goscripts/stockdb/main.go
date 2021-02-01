package main

import (
	"fmt"
	"os"
	"snack.com/xiyanxiyan10/conver"
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
	fmt.Printf("%s\n", util.Struct2Json(vec))
	for i, data := range vec {
		if i == 0 {
			continue
		}
		err = exchange.BackPutOHLC(conver.Int64Must(data[0]), conver.Float64Must(data[1]), conver.Float64Must(data[2]),
			conver.Float64Must(data[3]), conver.Float64Must(data[4]), conver.Float64Must(data[5]), "", exchange.GetPeriod())
		if err != nil {
			fmt.Printf("put ohlc fail:%s\n", err.Error())
			return
		}
	}
	/*
		exchange.SetSubscribe(symbol, constant.CacheAccount)
		exchange.SetSubscribe(symbol, constant.CacheRecord)
		exchange.SetSubscribe(symbol, constant.CachePosition)
		exchange.SetSubscribe(symbol, constant.CacheOrder)
		exchange.SetSubscribe(symbol, constant.CacheTicker)
	*/
}
