package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/trader"
)

func main() {
	var logger model.Logger
	var opt constant.Option
	var Constract = "quarter"
	var Symbol = "BTC/USD"
	//var Period = "M30"
	var IO = "online"
	var global trader.Global

	logger.Back = true
	global.Logger = logger
	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.LogBack = true

	maker := trader.ExchangeMaker[constant.FutureBack]
	exchange := maker(opt)

	exchange.SetIO(IO)
	exchange.SetContractType(Constract)
	exchange.SetStockType(Symbol)
	exchange.Ready()

	//global.Sleep(1000)
	err := exchange.BackGetStats()
	if err != nil {
		fmt.Printf("link to stockdb fail:%s", err.Error())
		return
	}
	markets, err := exchange.BackGetMarkets()
	if err != nil {
		fmt.Printf("fail to get markets:%s", err.Error())
		return
	}
	fmt.Printf("success to get markets:%v", markets)

	if len(markets) <= 0 {
		fmt.Printf("markets not found")
		return
	}
	symbols, err := exchange.BackGetSymbols()
	if err != nil {
		fmt.Printf("fail to get symbol:%s", err.Error())
		return
	}
	fmt.Printf("success to get symbol:%v", symbols)

	timeRange, err := exchange.BackGetTimeRange()
	if err != nil {
		fmt.Printf("fail to get timeRange:%s", err.Error())
		return
	}
	fmt.Printf("success to get timeRange:%v", timeRange)

	periodRange, err := exchange.BackGetPeriodRange()
	if err != nil {
		fmt.Printf("fail to get periodRange:%s", err.Error())
		return
	}
	fmt.Printf("success to get periodRange:%v", periodRange)
	return
}
