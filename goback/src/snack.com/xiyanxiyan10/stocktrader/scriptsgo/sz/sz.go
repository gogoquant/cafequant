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
	var Symbol = "sz000001"
	//var Period = "M30"
	var IO = "online"
	var global trader.Global

	logger.Back = true
	global.Logger = logger
	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.SZ
	opt.TraderID = 1
	opt.Type = constant.SZ
	opt.Index = 1
	opt.LogBack = true

	maker := trader.ExchangeMaker[opt.Type]
	exchange, err := maker(opt)
	if err != nil {
		fmt.Printf("init exchange fail:%s", err.Error())
		return
	}
	exchange.SetIO(IO)
	exchange.SetContractType(Constract)
	exchange.SetStockType(Symbol)
	exchange.Ready()

	records, err := exchange.GetRecords("M5", "")
	if err != nil {
		fmt.Printf("get records fail:%s", err.Error())
		return
	}
	fmt.Printf("success to get records:%v", records)

	ticker, err := exchange.GetTicker()
	if err != nil {
		fmt.Printf("get ticker fail:%s", err.Error())
		return
	}
	fmt.Printf("success to get ticker:%v", ticker)

	depth, err := exchange.GetDepth()
	if err != nil {
		fmt.Printf("get depth fail:%s", err.Error())
		return
	}
	fmt.Printf("success to get depth:%v", depth)
}
