package util

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
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

func (m *KlineMerge) Append(klines ...constant.Record) {
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
			}
			m.Count = 0
			m.curr = nil
		}

	}
}

func (m *KlineMerge) Get() []constant.Record {
	return m.vec
}
