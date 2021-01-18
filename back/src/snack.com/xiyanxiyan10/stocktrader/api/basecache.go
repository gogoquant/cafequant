package api

import (
	"sync"
	"time"

	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// BaseExchangeCache store the date from api in cahce
type BaseExchangeCache struct {
	Data      interface{}
	TimeStamp time.Time
	Mark      string
}

// BaseExchangeCaches ...
type BaseExchangeCaches struct {
	mutex    sync.Mutex
	depth    map[string]BaseExchangeCache
	position map[string]BaseExchangeCache
	account  map[string]BaseExchangeCache
	record   map[string]BaseExchangeCache
	order    map[string]BaseExchangeCache
	trader   map[string]BaseExchangeCache
	kline    map[string]BaseExchangeCache
	ticker   map[string]BaseExchangeCache
	//caches map[string]BaseExchangeCache
}

// Subscribe ...
func (e *BaseExchangeCaches) Subscribe() interface{} {
	return nil
}

// GetCache get ws val from cache
func (e *BaseExchangeCaches) GetCache(key string, stockSymbol string) BaseExchangeCache {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if key == constant.CacheTicker {
		return e.ticker[stockSymbol]
	}

	if key == constant.CachePosition {
		return e.position[stockSymbol]
	}

	if key == constant.CacheAccount {
		return e.account[""]
	}

	if key == constant.CacheRecord {
		return e.record[stockSymbol]
	}

	if key == constant.CacheOrder {
		return e.order[stockSymbol]
	}
	return BaseExchangeCache{}
}

// SetCache set ws val into cache
func (e *BaseExchangeCaches) SetCache(key string, stockSymbol string, val interface{}, mark string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	var item BaseExchangeCache

	item.Data = val
	item.TimeStamp = time.Now()
	item.Mark = mark

	if key == constant.CacheTicker {
		e.ticker[stockSymbol] = item
	}

	if key == constant.CachePosition {

		e.position[stockSymbol] = item
	}

	if key == constant.CacheAccount {
		e.account[""] = item
	}

	if key == constant.CacheRecord {
		e.record[stockSymbol] = item
	}

	if key == constant.CacheOrder {
		e.order[stockSymbol] = item
	}
}
