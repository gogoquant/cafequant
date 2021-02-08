package main

import (
	"fmt"
	"os"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"time"
)

// LoaderStragey ...
type LoaderStragey struct {
	Period    string
	Exchanges []api.Exchange
	Status    bool
}

// Init ...
func (e *LoaderStragey) Init(v map[string]string) error {
	period := v["period"]
	constract := v["constract"]
	symbol := v["symbol"]
	io := v["io"]
	exchange := e.Exchanges[0]
	exchange.SetIO(io)
	exchange.SetStockType(symbol + "." + constract)
	exchange.Start()
	exchange.SetLimit(1000)
	exchange.SetPeriod(period)
	exchange.SetPeriodSize(3)
	exchange.SetSubscribe(symbol, constant.CacheTicker)
	exchange.SetSubscribe(symbol, constant.CacheRecord)

	exchange.Log(constant.INFO, "", 0.0, 0.0, "Init success")
	e.Status = true
	return nil
}

// Run ...
func (e *LoaderStragey) Run(map[string]string) error {
	exchange := e.Exchanges[0]
	exchange.Log(constant.INFO, "", 0.0, 0.0, "Call")
	for e.Status {
		err := putOHLC(exchange, e.Period)
		if err != nil {
			exchange.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
	exchange.Log(constant.INFO, "", 0.0, 0.0, "Run stragey stop success")
	return nil
}

func putOHLC(exchange api.Exchange, period string) error {
	records, err := exchange.GetTicker()
	if err != nil {
		fmt.Printf("get records fail:%v\n", err)
		return err
	}
	fmt.Printf("get records success:%v\n", records)
	/*
		for _, record := range records {
			err = exchange.BackPutOHLC(record.Time, record.Open,
				record.High, record.Low, record.Close, record.Volume, "unknown", period)
			if err != nil {
				fmt.Printf("put ohlc to stockdb fail:%s", err.Error())
				return err
			}
		}
	*/
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
	var logger model.Logger
	var opt constant.Option
	var constract = "quarter"
	var symbol = "BTC/USD"
	var io = "block"
	var period = "M5"
	logger.Back = true

	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.BackLog = true
	opt.BackTest = true

	exchange, err := api.GetExchange(opt)
	if err != nil {
		fmt.Printf("init exchange fail:%s\n", err.Error())
		return
	}

	var loader LoaderStragey

	param := make(map[string]string)
	param["io"] = io
	param["symbol"] = symbol
	param["constract"] = constract
	param["period"] = period

	loader.Exchanges = append(loader.Exchanges, exchange)
	loader.Init(param)
	loader.Run(nil)
}
