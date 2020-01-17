package gobacktest

import (
	"errors"
	"github.com/xiyanxiyan10/quantcore/constant"
	"math"
)

// SizeHandler is the basic interface for setting the size of an order
type SizeHandler interface {
	SizeOrder(OrderEvent, DataEvent, PortfolioHandler) (*Order, error)
}

// Size is a basic size handler implementation
type Size struct {
	DefaultSize  int64
	DefaultValue float64
}

// SizeOrder adjusts the size of an order
func (s *Size) SizeOrder(order OrderEvent, data DataEvent, pf PortfolioHandler) (*Order, error) {
	// assert interface to concrete Type
	o := order.(*Order)
	price := order.Price()
	if price < 0 {
		o.SetOrderType(MarketOrder)
		price = data.Price()
	} else {
		o.SetOrderType(LimitOrder)
	}
	o.SetPrice(price)
	size := s.setDefaultSize(order.Qty(), price)
	// decide on order direction
	switch o.Direction() {
	case constant.TradeTypeLong:
		o.SetDirection(constant.TradeTypeLong)
		o.SetQty(float64(size))
	case constant.TradeTypeShort:
		o.SetDirection(constant.TradeTypeShort)
		o.SetQty(float64(size) * -1)
	default:
		return o, errors.New("unknown tradeType :" + string(o.Direction()))
	}
	return o, nil
}

// setDefaultSize ...
func (s *Size) setDefaultSize(size, price float64) int64 {
	correctedQty := int64(math.Floor(size / price))
	return correctedQty
}
