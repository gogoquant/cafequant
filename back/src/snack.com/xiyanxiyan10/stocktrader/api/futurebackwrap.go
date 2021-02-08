package api

import (
	"fmt"
	"snack.com/xiyanxiyan10/conver"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// NewFutureBackExchange create an exchange struct of futureExchange.com
func NewFutureBackExchange(opt constant.Option) (Exchange, error) {
	exchange := NewExchangeFutureBackLink(opt)
	if err := exchange.Init(opt); err != nil {
		return nil, err
	}
	return exchange, nil
}

// ExchangeFutureBackLink ...
type ExchangeFutureBackLink struct {
	ExchangeFutureBack

	records map[string][]constant.Record
}

// NewExchangeFutureBackLink create an exchange struct of futureExchange.com
func NewExchangeFutureBackLink(opt constant.Option) *ExchangeFutureBackLink {
	futureExchange := ExchangeFutureBackLink{
		records: make(map[string][]constant.Record),
	}
	opt.Limit = 10.0
	//futureExchange.BaseExchange.Init(opt)
	futureExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	futureExchange.SetID(opt.Index)
	return &futureExchange
}

// Init init the instance of this exchange
func (e *ExchangeFutureBackLink) Init(opt constant.Option) error {
	e.BaseExchange.Init(opt)
	return nil
}

// GetDepth get depth from exchange
func (e *ExchangeFutureBackLink) GetDepth() (*constant.Depth, error) {
	stockType := e.GetStockType()
	depth, err := e.ExchangeFutureBack.GetDepth(constant.DepthSize, stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is stockType")
		return nil, fmt.Errorf("GetDepth() error, the error number is stockType")
	}
	return depth, nil
}

// GetPosition get position from exchange
func (e *ExchangeFutureBackLink) GetPosition() ([]constant.Position, error) {
	return nil, fmt.Errorf("wait todo")
}

// GetMinAmount get the min trade amount of this exchange
func (e *ExchangeFutureBackLink) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *ExchangeFutureBackLink) GetAccount() (*constant.Account, error) {
	account, err := e.ExchangeFutureBack.GetAccount()
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetAccount() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetAccount() error, the error number is %s", err.Error())
	}
	return account, nil
}

// Buy buy from exchange
func (e *ExchangeFutureBackLink) Buy(price, amount string, msg string) (string, error) {
	var err error
	var ord *constant.Order
	stockType := e.GetStockType()
	if err := e.ValidBuy(); err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return "", fmt.Errorf("Buy() error, the error number is:%s", err.Error())
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
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return "", fmt.Errorf("Buy() error, the error number is:%s", err.Error())
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg)
	return ord.Id, nil

}

// Sell sell from exchange
func (e *ExchangeFutureBackLink) Sell(price, amount string, msg string) (string, error) {
	var err error
	var ord *constant.Order
	stockType := e.GetStockType()
	if err := e.ValidSell(); err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Sell() error, the error number is ", err.Error())
		return "", fmt.Errorf("Sell() error, the error number is:%s", err.Error())
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
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return "", fmt.Errorf("Buy() error, the error number is:%s", err.Error())
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg)
	return ord.Id, nil
}

// GetOrder get detail of an order
func (e *ExchangeFutureBackLink) GetOrder(id string) (*constant.Order, error) {
	order, err := e.ExchangeFutureBack.GetOneOrder(id, e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, conver.Float64Must(id), "GetOrder() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetOrder() error, the error number is:%s", err.Error())
	}
	return order, nil
}

// GetOrders get all unfilled orders
func (e *ExchangeFutureBackLink) GetOrders() ([]constant.Order, error) {
	orders, err := e.ExchangeFutureBack.GetUnfinishOrders(e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetOrders() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetOrders() error, the error number is:%s", err.Error())
	}
	return orders, nil
}

// CancelOrder cancel an order
func (e *ExchangeFutureBackLink) CancelOrder(orderID string) (bool, error) {
	result, err := e.ExchangeFutureBack.CancelOrder(orderID, e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, conver.Float64Must(orderID), "CancelOrder() error, the error number is ", err.Error())
		return false, fmt.Errorf("CancelOrder() error, the error number is:%s", err.Error())
	}
	if !result {
		return false, fmt.Errorf("CancelOrder() error, the error number is ")
	}
	e.logger.Log(constant.TradeTypeCancel, e.GetStockType(), 0, conver.Float64Must(orderID), "CancelOrder() success")
	return false, nil
}

// GetTicker get market ticker
func (e *ExchangeFutureBackLink) GetTicker() (*constant.Ticker, error) {
	ticker, err := e.ExchangeFutureBack.GetTicker(e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetTicker() error, the error number is %s", err.Error())
	}
	return ticker, nil
}

// GetRecords get candlestick data
func (e *ExchangeFutureBackLink) GetRecords() ([]constant.Record, error) {
	records, err := e.ExchangeFutureBack.GetRecords()
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			"GetRecords() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetRecords() error, the error number is %s", err.Error())
	}
	return records, nil
}
