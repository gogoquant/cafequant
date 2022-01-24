package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

// NewFutureBackExchange create an exchange struct of futureExchange.com
func NewFutureBackExchange(opt constant.Option) (Exchange, error) {
	exchange := NewExchangeFutureBackWrap(opt)
	if err := exchange.Init(opt); err != nil {
		return nil, err
	}
	return exchange, nil
}

// ExchangeFutureBackWrap ...
type ExchangeFutureBackWrap struct {
	ExchangeFutureBack
}

// NewExchangeFutureBackWrap create an exchange struct of futureExchange.com
func NewExchangeFutureBackWrap(opt constant.Option) *ExchangeFutureBackWrap {
	futureExchange := ExchangeFutureBackWrap{}
	futureExchange.currData = make(map[string]constant.OHLC)
	opt.Limit = 10.0
	futureExchange.BaseExchange.Init(opt)
	futureExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	futureExchange.SetID(opt.Index)
	return &futureExchange
}

// GetDepth get depth from exchange
func (e *ExchangeFutureBackWrap) GetDepth() (*constant.Depth, error) {
	stockType := e.GetStockType()
	depth, err := e.ExchangeFutureBack.GetDepth(constant.DepthSize, stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is stockType")
		return nil, nil
	}
	return depth, nil
}

// GetPosition get position from exchange
func (e *ExchangeFutureBackWrap) GetPosition() ([]constant.Position, error) {
	return nil, nil
}

// GetMinAmount get the min trade amount of this exchange
func (e *ExchangeFutureBackWrap) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *ExchangeFutureBackWrap) GetAccount() (*constant.Account, error) {
	account, err := e.ExchangeFutureBack.GetAccount()
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetAccount() error, the error number is ", err.Error())
		return nil, nil
	}
	return account, nil
}

// Buy buy from exchange
func (e *ExchangeFutureBackWrap) Buy(price, amount string, msg string) (string, error) {
	var err error
	var ord *constant.Order
	stockType := e.GetStockType()
	if err := e.ValidBuy(); err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(),
			util.Float64Must(amount), util.Float64Must(amount),
			"Buy() error, the error number is ", err.Error())
		return "", nil
	}
	if e.ExchangeFutureBack.GetDirection() == constant.TradeTypeLong {
		if price == "-1" {
			ord, err = e.ExchangeFutureBack.MarketBuy(amount, price, stockType)
		} else {
			ord, err = e.ExchangeFutureBack.LimitBuy(amount, price, stockType)
		}
	}

	if e.ExchangeFutureBack.GetDirection() == constant.TradeTypeShortClose {
		if price == "-1" {
			ord, err = e.ExchangeFutureBack.MarketSell(amount, price, stockType)
		} else {
			ord, err = e.ExchangeFutureBack.LimitSell(amount, price, stockType)
		}
	}

	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount),
			util.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return "", nil
	}
	priceFloat := util.Float64Must(price)
	amountFloat := util.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg)
	return ord.Id, nil

}

// Sell sell from exchange
func (e *ExchangeFutureBackWrap) Sell(price, amount string, msg string) (string, error) {
	var err error
	var ord *constant.Order
	stockType := e.GetStockType()
	if err := e.ValidSell(); err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount),
			util.Float64Must(amount), "Sell() error, the error number is ", err.Error())
		return "", nil
	}
	if e.ExchangeFutureBack.GetDirection() == constant.TradeTypeShort {
		if price == "-1" {
			ord, err = e.ExchangeFutureBack.MarketSell(amount, price, stockType)
		} else {
			ord, err = e.ExchangeFutureBack.LimitSell(amount, price, stockType)
		}
	}

	if e.ExchangeFutureBack.GetDirection() == constant.TradeTypeLongClose {
		if price == "-1" {
			ord, err = e.ExchangeFutureBack.MarketBuy(amount, price, stockType)
		} else {
			ord, err = e.ExchangeFutureBack.LimitBuy(amount, price, stockType)
		}
	}

	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount),
			util.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return "", nil
	}
	priceFloat := util.Float64Must(price)
	amountFloat := util.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg)
	return ord.Id, nil
}

// GetOrder get detail of an order
func (e *ExchangeFutureBackWrap) GetOrder(id string) (*constant.Order, error) {
	order, err := e.ExchangeFutureBack.GetOneOrder(id, e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, util.Float64Must(id),
			"GetOrder() error, the error number is ", err.Error())
		return nil, nil
	}
	return order, nil
}

// GetOrders get all unfilled orders
func (e *ExchangeFutureBackWrap) GetOrders() ([]constant.Order, error) {
	orders, err := e.ExchangeFutureBack.GetUnfinishOrders(e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			"GetOrders() error, the error number is ", err.Error())
		return nil, nil
	}
	return orders, nil
}

// CancelOrder cancel an order
func (e *ExchangeFutureBackWrap) CancelOrder(orderID string) (bool, error) {
	result, err := e.ExchangeFutureBack.CancelOrder(orderID, e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, util.Float64Must(orderID),
			"CancelOrder() error, the error number is ", err.Error())
		return false, nil
	}
	if !result {
		return false, nil
	}
	e.logger.Log(constant.TradeTypeCancel, e.GetStockType(), 0, util.Float64Must(orderID),
		"CancelOrder() success")
	return false, nil
}

// GetTicker get market ticker
func (e *ExchangeFutureBackWrap) GetTicker() (*constant.Ticker, error) {
	ticker, err := e.ExchangeFutureBack.GetTicker(e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ",
			err.Error())
		return nil, nil
	}
	return ticker, nil
}
