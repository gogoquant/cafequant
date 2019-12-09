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
	// no default set, no sizing possible, order rejected
	if (s.DefaultSize == 0) || (s.DefaultValue == 0) {
		return o, errors.New("cannot size order: no defaultSize or defaultValue set")
	}

	price := order.Price()
	if price < 0 {
		o.SetOrderType(MarketOrder)
		price = data.Price()
	} else {
		o.SetOrderType(LimitOrder)
	}

	// decide on order direction
	switch o.Direction() {
	case constant.TradeTypeLong:
		o.SetDirection(constant.TradeTypeLong)
		o.SetQty(order.Qty())
	case constant.TradeTypeShort:
		o.SetDirection(constant.TradeTypeShort)
		o.SetQty(order.Qty() * -1)
	default:
		return o, errors.New("unknown tradeType :" + string(o.Direction()))
	}
	return o, nil
}

func (s *Size) setDefaultSize(price float64) int64 {
	if (float64(s.DefaultSize) * price) > s.DefaultValue {
		correctedQty := int64(math.Floor(s.DefaultValue / price))
		return correctedQty
	}
	return s.DefaultSize
}
