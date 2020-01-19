package api

import (
	"github.com/miaolz123/stockdb/types"
)

// Driver is a stockdb interface
type Driver interface {
	close() error

	PutOHLC(datum types.OHLC, opt types.Option) response
	PutOHLCs(data []types.OHLC, opt types.Option) response
	PutOrder(datum types.Order, opt types.Option) response
	PutOrders(data []types.Order, opt types.Option) response
	GetStats() response
	GetMarkets() response
	GetSymbols(market string) response
	GetTimeRange(opt types.Option) response
	GetPeriodRange(opt types.Option) response
	GetOHLCs(opt types.Option) response
	GetDepth(opt types.Option) response
}
