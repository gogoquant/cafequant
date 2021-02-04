package main

import (
	"fmt"
	"os"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/goplugin"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/util"
	"time"
)

// TrendStragey ...
type TrendStragey struct {
	goplugin.GoStragey
	Status bool
}

// NewHandler ...
func NewHandler() (goplugin.GoStrageyHandler, error) {
	trend := new(TrendStragey)
	return trend, nil
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
	//exchange.Start()

	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Init success")
	e.Status = true
	return nil
}

// Run ...
func (e *TrendStragey) Run(map[string]string) error {
	exchange := e.Exchanges[0]
	exchange.Start()

	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Call")
	symbols, err := exchange.BackGetSymbols()
	if err != nil {
		e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Back get symbols fail")
		return nil
	}
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Back get symbols success:", symbols)

	times, err := exchange.BackGetTimeRange()
	if err != nil {
		e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Back get times fail")
		return nil
	}
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Back get times success:", times)

	str1 := time.Unix(times[0], 0).Local().String()
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Start time is:", str1)

	str := time.Unix(times[1], 0).Local().String()
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "End time is:", str)

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
		e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Display err is:", err.Error())
		return err
	}
	return nil
}

// Exit ...
func (e *TrendStragey) Exit(map[string]string) error {
	exchange := e.Exchanges[0]
	exchange.Stop()

	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Exit success")
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

	trend, err := NewHandler()
	if err != nil {
		fmt.Printf("create trend fail:%s\n", err.Error())
		return
	}
	param := make(map[string]string)
	param["io"] = io
	param["symbol"] = symbol
	param["constract"] = constract
	param["period"] = period
	trend.AddExchange(exchange)

	var logger model.Logger

	logger.Back = true

	trend.AddLogger(&logger)
	trend.Init(param)
	trend.Run(nil)
}

// Process ...
func (e *TrendStragey) Process() error {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Call")
	exchange := e.Exchanges[0]
	exchange.GetRecords()
	for e.Status {
		time.Sleep(time.Duration(3) * time.Minute)
	}
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Run stragey stop success")
	return nil
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
	exchange.Sleep(time.Minute * 1)
	positions, err := exchange.GetPosition()
	if err != nil {
		return fmt.Errorf("get position  fail:%s", err.Error())
	}
	if len(positions) == 0 {
		e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Try to close position and order")
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
