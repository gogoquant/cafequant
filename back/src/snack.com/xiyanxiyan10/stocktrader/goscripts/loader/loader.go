package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/api"
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

// NewLoaderHandler ...
func NewLoaderHandler(...interface{}) (goplugin.GoStrageyHandler, error) {
	loader := new(LoaderStragey)
	return loader, nil
}

// Init ...
func (e *LoaderStragey) Init(v ...interface{}) interface{} {
	if len(v) < 0 {
		e.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Parameter period needed")
		return nil
	}
	period, ok := v[0].(string)
	if !ok {
		e.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Parameter period convert fail")
		return nil
	}
	e.Period = period
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Init success")
	e.Status = true
	return "success"
}

// Run ...
func (e *LoaderStragey) Run(v ...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Call")
	exchange := e.Exchanges[0]
	for e.Status {
		err := putOHLC(exchange, e.Period)
		if err != nil {
			e.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			continue
		}
		time.Sleep(time.Duration(3) * time.Minute)
	}
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Run stragey stop success")
	return "success"
}

// Exit ...
func (e *LoaderStragey) Exit(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Exit success")
	e.Status = false
	return "success"
}

func putOHLC(exchange api.Exchange, period string) error {
	records, err := exchange.GetRecords(period, "", 3)
	if err != nil {
		fmt.Printf("get records fail:%v", err)
		return err
	}
	fmt.Printf("get records:%s\n", util.Struct2Json(records))
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
	var logger model.Logger
	var opt constant.Option
	var Constract = "quarter"
	var Symbol = "BTC/USD"
	var IO = "online"

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
	exchange.SetIO(IO)
	exchange.SetContractType(Constract)
	exchange.SetStockType(Symbol)
	exchange.Ready()

	loader, err := NewLoaderHandler()
	if err != nil {
		fmt.Printf("create loader fail:%s\n", err.Error())
		return
	}
	loader.AddExchange(exchange)
	loader.AddLogger(&logger)
	loader.Init("M5")
	loader.Run()
}
