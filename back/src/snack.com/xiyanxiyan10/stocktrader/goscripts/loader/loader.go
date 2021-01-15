package main

import (
	"fmt"
	"os"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/goplugin"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/util"
	"time"
)

// LoaderStragey ...
type LoaderStragey struct {
	goplugin.GoStragey
	Period string
	Status bool
}

// NewHandler ...
func NewHandler() (goplugin.GoStrageyHandler, error) {
	loader := new(LoaderStragey)
	return loader, nil
}

// Init ...
func (e *LoaderStragey) Init(v map[string]string) error {
	period := v["period"]
	constract := v["constract"]
	symbol := v["symbol"]
	io := v["io"]
	exchange := e.Exchanges[0]
	exchange.SetIO(io)
	exchange.SetContractType(constract)
	exchange.SetStockType(symbol)
	exchange.Ready()

	e.Period = period
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Init success")
	e.Status = true
	return nil
}

// Run ...
func (e *LoaderStragey) Run(map[string]string) error {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Call")
	exchange := e.Exchanges[0]
	for e.Status {
		err := putOHLC(exchange, e.Period)
		if err != nil {
			e.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			time.Sleep(time.Duration(3) * time.Minute)
			continue
		}
		time.Sleep(time.Duration(3) * time.Minute)
	}
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Run stragey stop success")
	return nil
}

// Exit ...
func (e *LoaderStragey) Exit(map[string]string) error {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Exit success")
	e.Status = false
	return nil
}

func putOHLC(exchange api.Exchange, period string) error {
	records, err := exchange.GetRecords(period, "", 3)
	if err != nil {
		fmt.Printf("get records fail:%v", err)
		return err
	}
	for _, record := range records {
		err = exchange.BackPutOHLC(record.Time, record.Open,
			record.High, record.Low, record.Close, record.Volume, "unknown", period)
		if err != nil {
			fmt.Printf("put ohlc to stockdb fail:%s", err.Error())
			return err
		}
	}
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
	var io = "online"
	var period = "M5"
	logger.Back = true

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

	loader, err := NewHandler()
	if err != nil {
		fmt.Printf("create loader fail:%s\n", err.Error())
		return
	}
	param := make(map[string]string)
	param["io"] = io
	param["symbol"] = symbol
	param["constract"] = constract
	param["period"] = period
	loader.AddExchange(exchange)
	loader.AddLogger(&logger)
	loader.Init(param)
	loader.Run(nil)
}
