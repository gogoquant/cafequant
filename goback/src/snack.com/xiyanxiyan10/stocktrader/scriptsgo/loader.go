package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/trader"
)

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

func main() {
	var logger model.Logger
	var opt constant.Option
	var Constract = "quarter"
	var Symbol = "BTC/USD"
	//var Period = "M30"
	var IO = "online"
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

	maker := trader.ExchangeMaker[opt.Type]
	exchange := maker(opt)

	exchange.SetIO(IO)
	exchange.SetContractType(Constract)
	exchange.SetStockType(Symbol)
	exchange.Ready()

	//global.Sleep(1000)
	err := exchange.BackGetStats()
	if err != nil {
		fmt.Printf("link to stockdb fail:%s", err.Error())
		return
	}
	records, err := exchange.GetRecords("M5", "")
	if err != nil {
		fmt.Printf("get records fail:%s", err.Error())
		return
	}
	fmt.Printf("success to get records:%v", records)
	err = putOHLC(exchange, "M5")
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
