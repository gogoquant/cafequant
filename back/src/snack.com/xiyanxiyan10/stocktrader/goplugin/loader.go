package goplugin

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/util"
	"time"
)

// LoaderStragey ...
type LoaderStragey struct {
	GoStragey
	Status bool
}

// NewLoaderHandler ...
func NewLoaderHandler(...interface{}) (GoStrageyHandler, error) {
	loader := new(LoaderStragey)
	return loader, nil
}

// Init ...
func (e *LoaderStragey) Init(...interface{}) interface{} {
	if e.Logger == nil {
		e.Logger.Log(constant.INFO, "", 0.0, 0.0, "logger is nil")
		return nil
	}
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Init")
	e.Status = true
	return nil
}

// Run ...
func (e *LoaderStragey) Run(v ...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Call")
	exchange := e.Exchanges[0]
	for e.Status {
		err := putOHLC(exchange, "M5")
		if err != nil {
			e.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			continue
		}
		time.Sleep(time.Duration(3) * time.Minute)
	}
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "stragey exit")
	return nil
}

// Exit ...
func (e *LoaderStragey) Exit(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Exit")
	e.Status = false
	return nil
}

func putOHLC(exchange api.Exchange, period string) error {
	records, err := exchange.GetRecords(period, "", 3)
	if err != nil {
		fmt.Printf("get records fail:%v", err)
		return err
	}
	fmt.Printf("get records:%s\n", util.Struct2Json(records))
	for _, record := range records {
		err = exchange.BackPutOHLC(record.Time, record.Open, record.High, record.Low, record.Close, record.Volume, "unknown", period)
		if err != nil {
			fmt.Printf("put ohlc to stockdb fail:%s", err.Error())
			return err
		}
	}
	return nil
}

// RunLoader ...
func RunLoader() {
	var logger model.Logger
	var opt constant.Option
	var Constract = "quarter"
	var Symbol = "BTC/USD"
	var period = "M5"
	var IO = "online"
	var interval = 5

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
	markets, err := exchange.BackGetStats()
	if err != nil {
		fmt.Printf("exchange get markets fail:%s\n", err.Error())
		return
	}
	fmt.Printf("markets %v", markets)
	/*
		for {
			records, err := exchange.GetRecords(period, "", 3)
			if err != nil {
				fmt.Printf("get records fail:%s\n", err.Error())
				time.Sleep(time.Duration(interval) * time.Minute)
				continue
			}
			fmt.Printf("records:%s\n", util.Struct2Json(records))
			time.Sleep(time.Duration(interval) * time.Minute)
		}
	*/
	for {
		err = putOHLC(exchange, period)
		if err != nil {
			fmt.Printf("put ohlcs fail:%s\n", err.Error())
			return
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
