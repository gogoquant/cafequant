package main

import (
	"fmt"
	"os"

	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
	//"time"
)

// main ...
func main() {
	if len(os.Args) < 2 {
		config.Init("./config.ini")
	} else {
		if err := config.Init(os.Args[1]); err != nil {
			fmt.Printf("config init error is %s\n", err.Error())
			return
		}
	}
	var opt constant.Option
	var symbol = "BTC/USDT.quater"
	var period = "M30"

	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.BackLog = true
	opt.BackTest = true
	opt.WatchList = append(opt.WatchList, "BTC/USDT.quater")

	exchange, err := api.GetExchange(opt)
	if err != nil {
		fmt.Printf("init exchange fail:%s\n", err.Error())
		return
	}

	exchange.SetStockType(symbol)
	exchange.SetPeriod(period)
	exchange.SetPeriodSize(3)
	exchange.Log(constant.INFO, "", 0.0, 0.0, "Call")
	err = exchange.Start()
	if err != nil {
		exchange.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
		return
	}
	for {
		ticker, err := exchange.GetTicker()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, err.Error())
			continue
		}
		if ticker == nil {
			break
		}
		fmt.Printf("ticker %s\n", util.Struct2Json(*ticker))
	}
}
