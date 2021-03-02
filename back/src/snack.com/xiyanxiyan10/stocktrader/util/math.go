package util

import (
	"github.com/markcheno/go-talib"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

type PriceType int

const (
	InClose PriceType = iota + 1
	InHigh
	InLow
	InOpen
)

type KlineMerge struct {
	size  int
	Limit int
	Count int
	curr  *constant.Record
	vec   []constant.Record
}

func NewKlineMerge(rate int) *KlineMerge {
	var merge KlineMerge
	merge.Limit = rate
	return &merge
}

func (m *KlineMerge) Append(klines ...constant.Record) int {
	var tot int = 0
	for _, kline := range klines {
		if m.curr == nil {
			m.curr = new(constant.Record)
			*(m.curr) = kline
		} else {
			m.curr.Volume += kline.Volume
			m.curr.Time = kline.Time
			m.curr.Close = kline.Close
			if m.curr.High < kline.High {
				m.curr.High = kline.High
			}
			if m.curr.Low > kline.Low {
				m.curr.Low = kline.Low
			}
		}
		m.Count++
		if m.Count >= m.Limit {
			// merge one kline and store it
			if m.curr != nil {
				m.vec = append(m.vec, *(m.curr))
				m.size += 1
				tot++
			}
			m.Count = 0
			m.curr = nil
		}

	}
	return tot
}

func (m *KlineMerge) Get(size int) []constant.Record {
	if size < 0 || size >= m.size {
		return m.vec
	}
	size = m.size - size
	return m.vec[size:]
}

func Ma(data []constant.Record, inTimePeriod int, maType talib.MaType, priceTy PriceType) []float64 {
	return talib.Ma(realData(data, priceTy), inTimePeriod, maType)
}

func Atr(data []constant.Record, inTimePeriod int) []float64 {
	var (
		inHigh  []float64
		inLow   []float64
		inClose []float64
	)

	for i := len(data) - 1; i >= 0; i-- {
		k := data[i]
		inHigh = append(inHigh, k.High)
		inLow = append(inLow, k.Low)
		inClose = append(inClose, k.Close)
	}

	return talib.Atr(inHigh, inLow, inClose, inTimePeriod)
}

func Macd(data []constant.Record, inFastPeriod int,
	inSlowPeriod int, inSignalPeriod int, priceTy PriceType) (DIF, DEA, MACD []float64) {
	var macd []float64
	dif, dea, hist := talib.Macd(realData(data, priceTy), inFastPeriod, inSlowPeriod, inSignalPeriod)
	for _, item := range hist {
		macd = append(macd, item*2)
	}
	return dif, dea, macd
}

func Boll(data []constant.Record, inTimePeriod int, deviation float64, priceTy PriceType) (up, middle, low []float64) {
	return talib.BBands(realData(data, priceTy), inTimePeriod, deviation, deviation, 0)
}

func realData(data []constant.Record, priceTy PriceType) []float64 {
	var inReal []float64
	for i := len(data) - 1; i >= 0; i-- {
		k := data[i]
		switch priceTy {
		case InClose:
			inReal = append(inReal, k.Close)
		case InHigh:
			inReal = append(inReal, k.High)
		case InLow:
			inReal = append(inReal, k.Low)
		case InOpen:
			inReal = append(inReal, k.Open)
		default:
			panic("please set ema type")
		}
	}
	return inReal
}
