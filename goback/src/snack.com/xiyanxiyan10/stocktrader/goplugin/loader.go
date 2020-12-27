package goplugin

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

// LoaderStragey ...
type LoaderStragey struct {
	GoStragey
}

// NewLoaderHandler ...
func NewLoaderHandler(...interface{}) (GoStrageyHandler, error) {
	var echo LoaderStragey
	return &echo, nil
}

// Set ...
func (e *LoaderStragey) Set(string, interface{}) interface{} {
	return nil
}

// Init ...
func (e *LoaderStragey) Init(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Init")
	return nil
}

// Call ...
func (e *LoaderStragey) Call(name string, v ...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Call")
	return nil
}

// Exit ...
func (e *LoaderStragey) Exit(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Exit")
	return nil
}

func putOHLC(exchange api.Exchange, period string) error {
	records, err := exchange.GetRecords(period, "")
	if err != nil {
		fmt.Printf("get records fail:%v", err)
		return err
	}
	fmt.Printf("get records:%v", records)
	for _, record := range records {
		err = exchange.BackPutOHLC(record.Time, record.Open, record.High, record.Low, record.Close, record.Volume, "", period)
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
		fmt.Printf("init exchange fail:%s", err.Error())
		return
	}
	exchange.SetIO(IO)
	exchange.SetContractType(Constract)
	exchange.SetStockType(Symbol)
	exchange.Ready()

	//global.Sleep(1000)
	err = exchange.BackGetStats()
	if err != nil {
		fmt.Printf("link to stockdb fail:%s", err.Error())
		return
	}

	fmt.Printf("link to stockdb success")
	return
	records, err := exchange.GetRecords(period, "")

	if err != nil {
		fmt.Printf("get records fail:%s", err.Error())
		return
	}
	fmt.Printf("success to get records:%v", records)
	err = putOHLC(exchange, period)
	if err != nil {
		fmt.Printf("put ohlcs fail:%s", err.Error())
		return
	}
	return
	markets, err := exchange.BackGetMarkets()
	if err != nil {
		fmt.Printf("fail to get markets:%s", err.Error())
		return
	}
	fmt.Printf("success to get markets:%v", markets)

	if len(markets) > 0 {
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
	}
	return
}
