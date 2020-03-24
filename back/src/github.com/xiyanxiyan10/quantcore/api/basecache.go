package api

import (
	"github.com/xiyanxiyan10/quantcore/constant"
	"sync"
	"time"
)

// BaseExchangeCache store val as cache for callback api
type BaseExchangeCache struct {
	BaseExchange

	sync.Mutex
	depth     constant.Depth
	depthTime time.Time

	orders     []constant.Order
	ordersTime time.Time

	traders     []constant.Trader
	tradersTime time.Time

	ticker     constant.Ticker
	tickerTime time.Time
}
