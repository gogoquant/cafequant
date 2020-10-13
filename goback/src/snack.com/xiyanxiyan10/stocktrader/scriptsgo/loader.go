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
	var Period = "M30"
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

	for {
		records, err := huobiExchange.GetRecords(Period, "")
		if err != nil {
			fmt.Printf("get records fail:%v", err)
		}
		fmt.Printf("get records:%v", records)
		global.Sleep(1000)
		err = huobiExchange.BackGetStats()
		if err != nil {
			fmt.Printf("link to stockdb fail:%s", err.Error())
			continue
		}
	}

}
