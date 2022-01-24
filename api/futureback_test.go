package api

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
	//"time"
)

// TestFutureBack ...
func TestFutureBack(t *testing.T) {
	config.Init("./config.ini")
	var opt constant.Option
	var symbol = "BTC/USDT.quater"

	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.FutureBack
	opt.TraderID = 1
	opt.Type = constant.FutureBack
	opt.Index = 1
	opt.BackLog = true
	opt.BackTest = true

	exchange, err := GetExchange(opt)
	if err != nil {
		log.Errorf("init exchange fail:%s", err.Error())
		return
	}

	exchange.SetStockType(symbol)
	exchange.Log(constant.INFO, "", 0.0, 0.0, "Call")
	err = exchange.Start()
	if err != nil {
		exchange.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
		return
	}
	for {
		ticker, err := exchange.GetTicker()
		if err != nil {
			exchange.Log(constant.INFO, "", 0.0, 0.0, err.Error())
			continue
		}
		if ticker == nil {
			break
		}
		exchange.Log(constant.INFO, "", 0.0, 0.0, "ticker"+util.Struct2Json(*ticker))
	}
}
