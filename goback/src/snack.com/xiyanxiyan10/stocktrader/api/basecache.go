package api

import (
	"sync"
	"time"
)

// BaseExchangeCache store the date from api in cahce
type BaseExchangeCache struct {
	Data      interface{}
	TimeStamp time.Time
}

// BaseExchangeCacheManager store val as cache for callback api
type BaseExchangeCachManager struct {
	BaseExchange
	sync.Mutex
	key    string // which trigger
	caches map[string]BaseExchangeCache
}
