package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// NewSZExchange create an exchange struct of futureExchange.com
func NewSZExchange(opt constant.Option) (Exchange, error) {
	exchange := NewSZSpotExchange(opt)
	exchange.SetRecordsPeriodMap(map[string]int64{
		"M5":  5,
		"M15": 15,
		"M30": 30,
	})
	if err := exchange.Init(opt); err != nil {
		return nil, err
	}
	return exchange, nil
}
