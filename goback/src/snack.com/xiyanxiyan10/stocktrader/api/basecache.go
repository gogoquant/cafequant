package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"sync"
	"time"
)

// BaseExchangeCache store the date from api in cahce
type BaseExchangeCache struct {
	Data      interface{}
	TimeStamp time.Time
	Mark      string
}

// BaseExchangeCachePool ...
type BaseExchangeCachePool struct {
	mutex  sync.Mutex
	depth  map[string]BaseExchangeCache
	order  map[string]BaseExchangeCache
	trader map[string]BaseExchangeCache
	kline  map[string]BaseExchangeCache
	ticker map[string]BaseExchangeCache
	//caches map[string]BaseExchangeCache
}

// Subscribe ...
func (e *BaseExchangeCachePool) Subscribe() interface{} {
	return nil
}

// GetCache get ws val from cache
func (e *BaseExchangeCachePool) GetCache(key string, stockSymbol string) BaseExchangeCache {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if key == constant.CacheDepth {
		return e.depth[stockSymbol]
	}

	if key == constant.CacheTicker {
		return e.ticker[stockSymbol]
	}

	if key == constant.CacheTrader {
		return e.trader[stockSymbol]
	}

	if key == constant.CacheKline {
		return e.kline[stockSymbol]
	}
	return BaseExchangeCache{}
}

// SetCache set ws val into cache
func (e *BaseExchangeCachePool) SetCache(key string, stockSymbol string, val interface{}, mark string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	var item BaseExchangeCache

	item.Data = val
	item.TimeStamp = time.Now()
	item.Mark = mark

	if key == constant.CacheDepth {
		e.depth[stockSymbol] = item
	}

	if key == constant.CacheTicker {
		e.ticker[stockSymbol] = item
	}

	if key == constant.CacheTrader {
		e.trader[stockSymbol] = item
	}

	if key == constant.CacheKline {
		e.kline[stockSymbol] = item
	}
}
