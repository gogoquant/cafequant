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
