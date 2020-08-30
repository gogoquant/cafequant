package api

import (
	"github.com/nntaoli-project/goex"
	"time"
)

type DataConfig struct {
	Ex       string
	Pair     goex.CurrencyPair
	StarTime time.Time
	EndTime  time.Time
	Size     int //多少档深度数据
	UnGzip   bool
}
