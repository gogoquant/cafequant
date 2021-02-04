package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"sync"
	"time"
)

// BaseExchangeCache store the date from api in cahce
type BaseExchangeCache struct {
	Data      interface{}
	TimeStamp int64
	Mark      bool
}

// BaseExchangeCaches ...
type BaseExchangeCaches struct {
	mutex      sync.Mutex
	ch         chan [2]string
	waitsymbol string
	waitaction string

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

// Wait set wait symbol + action
func (e *BaseExchangeCaches) wait(symbol, action string) {
	e.waitsymbol = symbol
	e.waitaction = action
}

// PUSH push val into
func (e *BaseExchangeCaches) push(symbol, action string) {
	if e.waitsymbol == symbol && e.waitaction == action {
		e.ch <- [2]string{symbol, action}
	}
	return
}

// POP val out
func (e *BaseExchangeCaches) pop(symbol, action string) interface{} {
	val := <-e.ch
	return val
}

// GetCache get ws val from cache
func (e *BaseExchangeCaches) GetCache(action string, stockSymbol string, fresh bool) BaseExchangeCache {
	var dst BaseExchangeCache
	if fresh {
		// block wait
		for {
			val := e.pop(stockSymbol, action)
			if val == nil {
				// the chan close
				dst.Data = nil
				return dst
			}
			vec := val.([2]string)
			if vec[0] != stockSymbol || vec[1] != action {
				//wait the next one
				continue
			}
			break
		}
	}

	e.mutex.Lock()

	if action == constant.CacheTicker {
		dst = e.ticker[stockSymbol]
	}

	if action == constant.CachePosition {
		dst = e.position[stockSymbol]
	}

	if action == constant.CacheAccount {
		dst = e.account[""]
	}

	if action == constant.CacheRecord {
		dst = e.record[stockSymbol]
	}

	if action == constant.CacheOrder {
		dst = e.order[stockSymbol]
	}
	if !dst.Mark {
		dst.Data = nil
	}
	e.mutex.Unlock()

	return dst
}

// SetCache set ws val into cache
func (e *BaseExchangeCaches) SetCache(action string, stockSymbol string, val interface{}, fresh bool) {
	//lock
	e.mutex.Lock()
	var item BaseExchangeCache

	item.Data = val
	item.TimeStamp = time.Now().Unix()
	item.Mark = true

	if action == constant.CacheTicker {
		if e.ticker == nil {
			e.ticker = make(map[string]BaseExchangeCache)
		}
		e.ticker[stockSymbol] = item
	}

	if action == constant.CachePosition {
		if e.position == nil {
			e.position = make(map[string]BaseExchangeCache)
		}
		e.position[stockSymbol] = item
	}

	if action == constant.CacheAccount {
		if e.account == nil {
			e.account = make(map[string]BaseExchangeCache)
		}
		e.account[""] = item
	}

	if action == constant.CacheRecord {
		if e.record == nil {
			e.record = make(map[string]BaseExchangeCache)
		}
		e.record[stockSymbol] = item
	}

	if action == constant.CacheOrder {
		if e.order == nil {
			e.order = make(map[string]BaseExchangeCache)
		}
		e.order[stockSymbol] = item
	}
	e.mutex.Unlock()
	// unlock

	if fresh {
		e.push(stockSymbol, action)
	}
}
