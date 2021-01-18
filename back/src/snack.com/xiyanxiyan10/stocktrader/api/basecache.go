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

	var dst BaseExchangeCache
	if key == constant.CacheTicker {
		dst = e.ticker[stockSymbol]
	}

	if key == constant.CachePosition {
		dst = e.position[stockSymbol]
	}

	if key == constant.CacheAccount {
		dst = e.account[""]
	}

	if key == constant.CacheRecord {
		dst = e.record[stockSymbol]
	}

	if key == constant.CacheOrder {
		dst = e.order[stockSymbol]
	}
	if len(dst.Mark) == 0 {
		dst.Data = nil
	}
	return dst
}

// SetCache set ws val into cache
func (e *BaseExchangeCaches) SetCache(key string, stockSymbol string, val interface{}, mark string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	var item BaseExchangeCache

	item.Data = val
	item.TimeStamp = time.Now()
	item.Mark = "mark"

	if key == constant.CacheTicker {
		if e.ticker == nil {
			e.ticker = make(map[string]BaseExchangeCache)
		}
		e.ticker[stockSymbol] = item
	}

	if key == constant.CachePosition {
		if e.position == nil {
			e.position = make(map[string]BaseExchangeCache)
		}
		e.position[stockSymbol] = item
	}

	if key == constant.CacheAccount {
		if e.account == nil {
			e.account = make(map[string]BaseExchangeCache)
		}
		e.account[""] = item
	}

	if key == constant.CacheRecord {
		if e.record == nil {
			e.record = make(map[string]BaseExchangeCache)
		}
		e.record[stockSymbol] = item
	}

	if key == constant.CacheOrder {
		if e.order == nil {
			e.order = make(map[string]BaseExchangeCache)
		}
		e.order[stockSymbol] = item
	}
}
