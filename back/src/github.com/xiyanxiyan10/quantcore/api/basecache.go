package api

import (
	"github.com/xiyanxiyan10/quantcore/constant"
	"sync"
)

// BaseExchangeCache store val as cache for callback api
type BaseExchangeCache struct {
	BaseExchange

	sync.Mutex
	depth   constant.Depth
	orders  []constant.Order
	traders []constant.Trader
	ticker  constant.Ticker
}
