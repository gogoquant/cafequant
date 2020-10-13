package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/trader"
)

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
	huobiExchange := maker(opt)

	huobiExchange.SetIO(IO)
	huobiExchange.SetContractType(Constract)
	huobiExchange.SetStockType(Symbol)
	huobiExchange.Ready()

	//global.Sleep(1000)
	err := huobiExchange.BackGetStats()
	if err != nil {
		fmt.Printf("link to stockdb fail:%s", err.Error())
		return
	}
	/*
		records, err := huobiExchange.GetRecords(Period, "")
		if err != nil {
			fmt.Printf("get records fail:%v", err)
			return
		}
		fmt.Printf("get records:%v", records)
			for _, record := range records {
				err = huobiExchange.BackPutOHLC(record.Time, record.Open, record.High, record.Low, record.Close, record.Volume, "", Period)
				if err != nil {
					fmt.Printf("put ohlc to stockdb fail:%s", err.Error())
					return
				}
			}
	*/
	markets, err := huobiExchange.BackGetMarkets()
	if err != nil {
		fmt.Printf("fail to get markets:%s", err.Error())
		return
	}
	fmt.Printf("success to get markets:%v", markets)

	if len(markets) > 0 {
		symbols, err := huobiExchange.BackGetSymbols()
		if err != nil {
			fmt.Printf("fail to get symbol:%s", err.Error())
			return
		}
		fmt.Printf("success to get symbol:%v", symbols)

		timeRange, err := huobiExchange.BackGetTimeRange()
		if err != nil {
			fmt.Printf("fail to get timeRange:%s", err.Error())
			return
		}
		fmt.Printf("success to get timeRange:%v", timeRange)

		periodRange, err := huobiExchange.BackGetPeriodRange()
		if err != nil {
			fmt.Printf("fail to get periodRange:%s", err.Error())
			return
		}
		fmt.Printf("success to get periodRange:%v", periodRange)
	}
	return
}
