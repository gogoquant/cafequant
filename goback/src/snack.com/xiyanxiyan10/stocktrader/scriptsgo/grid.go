package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/trader"
)

// main ...
func main() {
	var logger model.Logger
	var opt constant.Option
	var Constract = "quarter"
	var Symbol = "BTC/USD"
	var Period = "M30"
	var IO = "online"
	var Coin = "BTC"
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
	//exchange.Ready()

	//global.Sleep(1000)
	err = exchange.BackGetStats()
	if err != nil {
		fmt.Printf("link to stockdb fail:%s", err.Error())
		return
	}
	markets, err := exchange.BackGetMarkets()
	if err != nil {
		fmt.Printf("fail to get markets:%s", err.Error())
		return
	}
	fmt.Printf("success to get markets:%v\n", markets)

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
	exchange.SetBackTime(timeRange[0], timeRange[1], Period)
	exchange.SetBackCommission(0, 0, 100, 0.001, true)
	exchange.SetBackAccount(Coin, 1.0)
	exchange.SetMarginLevel(5)
	err = exchange.Ready()
	if err != nil {
		fmt.Printf("fail to back ready:%s", err.Error())
		return
	}
	for i := 0; i < 100; i++ {
		ticker, err := exchange.GetTicker()
		if err != nil {
			fmt.Printf("fail to get ticker:%s", err.Error())
			return
		}
		if ticker == nil {
			break
		}
		fmt.Printf("get ticker:%v\n", ticker)
		exchange.SetDirection("sell")
		ID, err := exchange.Sell("9000", "10", "buy")
		if err != nil {
			fmt.Printf("fail to buy ticker:%s", err.Error())
		} else {
			fmt.Printf("buy order put success\n" + ID)
			/*
				_, err := exchange.CancelOrder(ID)
				if err != nil {
					fmt.Printf("cancel order fail:%s", err.Error())
				}
			*/
		}
	}

	for i := 0; i < 2; i++ {
		ticker, err := exchange.GetTicker()
		if err != nil {
			fmt.Printf("fail to get ticker:%s\n", err.Error())
			return
		}
		if ticker == nil {
			break
		}
		exchange.SetDirection("closebuy")
		ID, err := exchange.Sell("11000", "5", "closebuy")
		if err != nil {
			fmt.Printf("fail to closebuy ticker:%s\n", err.Error())
		} else {
			fmt.Printf("closebuy order put success\n" + ID)
			/*
				_, err := exchange.CancelOrder(ID)
				if err != nil {
					fmt.Printf("cancel order fail:%s\n", err.Error())

				}
			*/
		}
	}
	for i := 0; i < 1000; i++ {
		ticker, err := exchange.GetTicker()
		if err != nil {
			fmt.Printf("fail to get ticker:%s\n", err.Error())
			return
		}
		if ticker == nil {
			break
		}

	}
	return
}
