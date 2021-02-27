package main

import (
	"fmt"
	"os"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
	"time"
)

// TrendStragey ...
type TrendStragey struct {
	Exchanges []api.Exchange
	Global    api.GlobalHandler
	Status    bool
}

// Init ...
func (e *TrendStragey) Init(v map[string]string, opt constant.Option) error {
	period := v["period"]
	constract := v["constract"]
	symbol := v["symbol"]
	io := v["io"]

	exchange := e.Exchanges[0]
	exchange.SetIO(io)
	key := symbol + "." + constract
	exchange.SetStockType(symbol + "." + constract)
	exchange.SetPeriod(period)
	exchange.SetPeriodSize(10)
	exchange.SetPeriodSize(5)
	exchange.SetSubscribe(key, constant.CacheAccount)
	exchange.SetSubscribe(key, constant.CacheRecord)
	exchange.SetSubscribe(key, constant.CachePosition)
	exchange.SetSubscribe(key, constant.CacheOrder)
	exchange.SetSubscribe(key, constant.CacheTicker)

	e.Global = api.NewGlobal(opt)

	exchange.Log(constant.INFO, "", 0.0, 0.0, "Init success")
	e.Status = true
	return nil
}

// Run ...
func (e *TrendStragey) Run() error {
	exchange := e.Exchanges[0]
	global := e.Global

	exchange.Log(constant.INFO, "", 0.0, 0.0, "Call")
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

	startStr := time.Unix(times[0], 0).Local().String()
	exchange.Log(constant.INFO, "", 0.0, 0.0, fmt.Sprintf("Start time is:%s", startStr))

	endStr := time.Unix(times[1], 0).Local().String()
	exchange.Log(constant.INFO, "", 0.0, 0.0, fmt.Sprintf("End time is:%s", endStr))

	exchange.SetBackTime(times[0], times[1], exchange.GetPeriod())

	global.DrawSetPath("./trend.html")
	exchange.Start()
	fmt.Printf("trend start\n")

	for {
		records, err := exchange.GetRecords()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, err.Error())
			break
		}
		if records == nil {
			fmt.Printf("get ohlc null\n")
			break
		}
		ohlcs := records
		for _, ohlc := range ohlcs {
			fmt.Printf("get ohlc %s\n", util.Struct2Json(ohlc))
			global.DrawKLine(util.TimeUnix2Str(ohlc.Time),
				float32(ohlc.Open), float32(ohlc.Close), float32(ohlc.Low), float32(ohlc.High))
		}

		err = global.DrawPlot()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, fmt.Sprintf("Display err is:%s", err.Error()))
			return err
		}
	}

	return nil
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
	var constract = "quarter"
	var symbol = "BTC/USD"
	var io = "online"
	var period = "H1"

	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.BackLog = true
	opt.BackTest = false

	exchange, err := api.GetExchange(opt)
	if err != nil {
		fmt.Printf("init exchange fail:%s\n", err.Error())
		return
	}

	var trend TrendStragey
	trend.Exchanges = append(trend.Exchanges, exchange)

	param := make(map[string]string)
	param["io"] = io
	param["symbol"] = symbol
	param["constract"] = constract
	param["period"] = period

	trend.Init(param, opt)
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
