package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/trader"
)

const (
	Constract = "quarter"
	Symbol    = "BTC/USD"
	Period    = "M30"
	IO        = "online"
	Coin      = "BTC"
)

func main() {
	var logger model.Logger
	var opt constant.Option
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
	exchange, err := maker(opt)
	if err != nil {
		fmt.Printf("init exchange fail:%s", err.Error())
		return
	}
	exchange.SetIO(IO)
	exchange.SetContractType(Constract)
	exchange.SetStockType(Symbol)
	exchange.SetSubscribe(Symbol, "ticker")

	timeRange, err := exchange.BackGetTimeRange()
	if err != nil {
		fmt.Printf("get time range fail:%s", err.Error())
		return
	}

	exchange.SetBackTime(timeRange[0], timeRange[1], Period)
	exchange.SetBackCommission(0, 0, 100, 0.001, true)
	exchange.SetBackAccount(Coin, 1.0)
	exchange.SetMarginLevel(5)

	err = exchange.Ready()
	if err != nil {
		fmt.Printf("fail to back ready:%s", err.Error())
		return
	}
}
