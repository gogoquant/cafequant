package api

import (
	goex "github.com/nntaoli-project/goex"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// NewHuoBiExchange create an exchange struct of futureExchange.com
func NewHuoBiExchange(opt constant.Option) (Exchange, error) {
	exchange := NewSpotExchange(opt)
	exchange.SetRecordsPeriodMap(map[string]int64{
		"M1":  goex.KLINE_PERIOD_1MIN,
		"M5":  goex.KLINE_PERIOD_5MIN,
		"M15": goex.KLINE_PERIOD_15MIN,
		"M30": goex.KLINE_PERIOD_30MIN,
		"H1":  goex.KLINE_PERIOD_1H,
		"H2":  goex.KLINE_PERIOD_4H,
		"H4":  goex.KLINE_PERIOD_4H,
		"D1":  goex.KLINE_PERIOD_1DAY,
		"W1":  goex.KLINE_PERIOD_1WEEK,
	})
	if err := exchange.Init(opt); err != nil {
		return nil, err
	}
	return exchange, nil
}
