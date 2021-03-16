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

// TrendStragey ...
type TrendStragey struct {
	Exchanges []api.Exchange
	Global    api.GlobalHandler

	Ticker  constant.Ticker
	Records []constant.Record

	Status bool
}

// Run ...
func (e *TrendStragey) Run() error {
	exchange := e.Exchanges[0]

	exchange.Log(constant.INFO, "", 0.0, 0.0, "Call")
	/*
		symbols, err := exchange.BackGetSymbols()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, "Back get symbols fail")
			return nil
		}
		exchange.Log(constant.INFO, "", 0.0, 0.0, fmt.Sprintf("Back get symbols success:%s", symbols))

		times, err := exchange.BackGetTimeRange()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, "Back get times fail")
			return nil
		}

		exchange.SetBackTime(times[0], times[1], exchange.GetPeriod())
	*/
	exchange.Start()

	fmt.Printf("trend start\n")

	for {
		records, err := exchange.GetRecords()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, err.Error())
			continue
		}
		if len(records) == 0 {
			exchange.Log(constant.INFO, "", 0.0, 0.0, "records not found\n")
			continue
		}

		fmt.Printf("record %s\n", util.Struct2Json(records[0]))
		ticker, err := exchange.GetTicker()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, err.Error())
			continue
		}
		if ticker == nil {
			continue
		}
		fmt.Printf("ticker %s\n", util.Struct2Json(*ticker))
	}
}

// Exit ...
func (e *TrendStragey) Exit(map[string]string) error {
	exchange := e.Exchanges[0]
	exchange.Stop()
	exchange.Log(constant.INFO, "", 0.0, 0.0, "Exit success")
	e.Status = false
	return nil
}

// main ...
func main() {
	if len(os.Args) < 2 {
		fmt.Println("命令行的参数不足:", len(os.Args))
		config.Init("./config.ini")
	} else {
		if err := config.Init(os.Args[1]); err != nil {
			fmt.Printf("config init error is %s\n", err.Error())
			return
		}
	}
	var opt constant.Option
	var symbol = "BTC/USD.quarter"
	var io = constant.IOBLOCK
	var period = "M30"

	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.BackLog = true
	opt.BackTest = false
	opt.BackExit = true

	exchange, err := api.GetExchange(opt)
	if err != nil {
		fmt.Printf("init exchange fail:%s\n", err.Error())
		return
	}

	var trend TrendStragey
	trend.Exchanges = append(trend.Exchanges, exchange)

	exchange.SetIO(io)

	exchange.SetStockType(symbol)
	exchange.SetPeriod(period)
	exchange.SetPeriodSize(3)
	//exchange.SetSubscribe(symbol, constant.CacheAccount)
	exchange.SetSubscribe(symbol, constant.CacheRecord)
	//exchange.SetSubscribe(symbol, constant.CachePosition)
	//exchange.SetSubscribe(symbol, constant.CacheOrder)
	exchange.SetSubscribe(symbol, constant.CacheTicker)
	trend.Global = api.NewGlobal(opt)

	exchange.Log(constant.INFO, "", 0.0, 0.0, "Init success")
	trend.Status = true
	trend.Run()
}

// Action trendAction put order and watch this order
func (e *TrendStragey) Action(low, high, amount float64, dir int) error {
	exchange := e.Exchanges[0]
	direction := "buy"
	closedirection := "closebuy"
	openFunc := exchange.Buy
	closeFunc := exchange.Sell
	openPrice := high
	closePrice := low
	if dir == 1 {
		direction = "sell"
		closedirection = "closesell"
		openFunc = exchange.Sell
		closeFunc = exchange.Buy
		openPrice = low
		closePrice = high
	}
	exchange.SetDirection(direction)
	_, err := openFunc(fmt.Sprintf("%f", openPrice), fmt.Sprintf("%f", amount), "open order")
	if err != nil {
		return fmt.Errorf("open order fail:%s", err.Error())
	}
	exchange.Sleep(1)
	positions, err := exchange.GetPosition()
	if err != nil {
		return fmt.Errorf("get position  fail:%s", err.Error())
	}
	if len(positions) == 0 {
		exchange.Log(constant.INFO, "", 0.0, 0.0, "Try to close position and order")
		// Todo close position here
	}
	exchange.SetDirection(closedirection)
	_, err = closeFunc(fmt.Sprintf("%f", closePrice), fmt.Sprintf("%f", amount), "open order")
	if err != nil {
		return fmt.Errorf("open order fail:%s", err.Error())
	}
	//exchange.Buy()
	return nil
}
