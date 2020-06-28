package api

import (
	"sync"
	"time"
)

const (
	// CacheTicker ...
	CacheTicker = "ticker"

	// CacheDepth ...
	CacheDepth = "depth"

	// CacheTrader ...
	CacheTrader = "trader"

	// CacheKline ...
	CacheKline = "kline"
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
	caches map[string]BaseExchangeCache
}

// Subscribe ...
func (e *BaseExchangeCachePool) Subscribe() interface{} {
	return nil
}

// Get get ws val from cache
func (e *BaseExchangeCachePool) GetCache(key string) BaseExchangeCache {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.caches[key]
}

// Set set ws val into cache
func (e *BaseExchangeCachePool) SetCache(key string, val interface{}, mark string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	var item BaseExchangeCache
	item.Data = val
	item.TimeStamp = time.Now()
	item.Mark = mark
	e.caches[key] = item
}
