package main

import (
	"fmt"
	"os"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/util"
	"time"
)

// TrendStragey ...
type TrendStragey struct {
	Exchanges []api.Exchange
	Status    bool
}

// Init ...
func (e *TrendStragey) Init(v map[string]string) error {
	period := v["period"]
	constract := v["constract"]
	symbol := v["symbol"]
	io := v["io"]

	exchange := e.Exchanges[0]
	exchange.SetIO(io)
	exchange.SetContractType(constract)
	exchange.SetStockType(symbol)
	exchange.SetPeriod(period)
	exchange.SetPeriodSize(10)
	exchange.SetSubscribe(symbol, constant.CacheAccount)
	exchange.SetSubscribe(symbol, constant.CacheRecord)
	exchange.SetSubscribe(symbol, constant.CachePosition)
	exchange.SetSubscribe(symbol, constant.CacheOrder)
	exchange.SetSubscribe(symbol, constant.CacheTicker)

	exchange.Log(constant.INFO, "", 0.0, 0.0, "Init success")
	e.Status = true
	return nil
}

// Run ...
func (e *TrendStragey) Run() error {
	exchange := e.Exchanges[0]
	exchange.Start()

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

	ohlcs, err := exchange.BackGetOHLCs(times[0], times[1], exchange.GetPeriod())
	if err != nil {
		return err
	}

	drawHandler := draw.NewDrawHandler()
	drawHandler.SetPath("/Users/shu/Desktop/trend.html")
	for i, ohlc := range ohlcs {
		drawHandler.PlotKLine(util.TimeUnix2Str(ohlc.Time),
			float32(ohlc.Open), float32(ohlc.Close), float32(ohlc.Low), float32(ohlc.High))
		drawHandler.PlotLine("low", util.TimeUnix2Str(ohlc.Time), float32(ohlc.Low), "")
		//drawHandler.PlotLine("high", util.TimeUnix2Str(ohlc.Time), float32(ohlc.High), "")
		if i > 1 && ohlcs[i].Low > ohlcs[i-1].Low {
			drawHandler.PlotLine("vol", util.TimeUnix2Str(ohlc.Time), 30000, draw.StepLine)
		} else {
			drawHandler.PlotLine("vol", util.TimeUnix2Str(ohlc.Time), 0, draw.StepLine)
		}
	}
	err = drawHandler.Display()
	if err != nil {
		exchange.Log(constant.INFO, "", 0.0, 0.0, fmt.Sprintf("Display err is:%s", err.Error()))
		return err
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
	var period = "M15"

	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.BackLog = true
	opt.BackTest = true

	maker := api.ExchangeMaker[opt.Type]
	exchange, err := maker(opt)
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
	trend.Init(param)
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
