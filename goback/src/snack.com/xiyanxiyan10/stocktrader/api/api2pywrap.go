package api

import (
	"errors"
	"github.com/qiniu/py"
	"github.com/qiniu/py/pyutil"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// BackTime ...
type BackTime struct {
	Start  int64
	End    int64
	Period string
}

// BackCommission ...
type BackCommission struct {
	Taker        float64
	Maker        float64
	ContractRate float64 // 合约每张价值
	Coin         bool
	CoverRate    float64 // 合约每张价值
}

// ExchangePythonLink ...
type ExchangePythonLink struct {
	api Exchange
}

// NewExchangePython create an exchange struct of futureExchange.com
func NewExchangePython(e func(opt constant.Option) (Exchange, error)) func(opt constant.Option) (ExchangePython, error) {
	return func(opt constant.Option) (ExchangePython, error) {
		var ex ExchangePythonLink
		api, err := e(opt)
		if err != nil {
			return nil, err
		}
		ex.api = api
		return &ex, nil
	}
}

// Ready ...
func (e *ExchangePythonLink) Ready(args *py.Tuple) (ret *py.Base, err error) {
	err = e.api.Ready()
	if err != nil {
		return py.IncNone(), nil
	}
	return py.IncNone(), nil
}

// SetIO ...
func (e *ExchangePythonLink) SetIO(args *py.Tuple) (ret *py.Base, err error) {
	var i string
	err = py.ParseV(args, &i)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetIO(i)
	return py.IncNone(), nil
}

// GetIO ...
func (e *ExchangePythonLink) GetIO(args *py.Tuple) (ret *py.Base, err error) {
	i := e.api.GetIO()
	val, ok := pyutil.NewVar(i)
	if !ok {
		return py.IncNone(), nil
	}
	return val, nil
}

// Subscribe ...
func (e *ExchangePythonLink) Subscribe(args *py.Tuple) (ret *py.Base, err error) {
	var stock, action string
	err = py.ParseV(args, &stock, &action)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetSubscribe(stock, action)
	return py.IncNone(), nil
}

// Log ...
func (e *ExchangePythonLink) Log(args *py.Tuple) (ret *py.Base, err error) {
	var msgs []interface{}
	err = py.ParseV(args, &msgs)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.Log(msgs...)
	return py.IncNone(), nil
}

// GetType ...
func (e *ExchangePythonLink) GetType(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetType()
	val, ok := pyutil.NewVar(s)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetName ...
func (e *ExchangePythonLink) GetName(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetType()
	val, ok := pyutil.NewVar(s)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// SetLimit ...
func (e *ExchangePythonLink) SetLimit(args *py.Tuple) (ret *py.Base, err error) {
	var vars []interface{}
	err = py.ParseV(args, &vars)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetLimit(vars)
	return py.IncNone(), nil
}

// Sleep ...
func (e *ExchangePythonLink) Sleep(args *py.Tuple) (ret *py.Base, err error) {
	var vars []interface{}
	err = py.ParseV(args, &vars)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.Sleep(vars)
	return py.IncNone(), nil
}

// GetAccount ...
func (e *ExchangePythonLink) GetAccount(args *py.Tuple) (ret *py.Base, err error) {
	account, err := e.api.GetAccount()
	if err != nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(*account)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetDepth ...
func (e *ExchangePythonLink) GetDepth(args *py.Tuple) (ret *py.Base, err error) {
	depth, err := e.api.GetDepth()
	if err == nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(*depth)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// Buy ...
func (e *ExchangePythonLink) Buy(args *py.Tuple) (ret *py.Base, err error) {
	var price, amount string
	var msgs []interface{}
	err = py.ParseV(args, &price, &amount, &msgs)
	if err != nil {
		return py.IncNone(), err
	}
	ID, err := e.api.Buy(price, amount, msgs)
	if err != nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(ID)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// Sell ...
func (e *ExchangePythonLink) Sell(args *py.Tuple) (ret *py.Base, err error) {
	var price, amount string
	var msgs []interface{}
	err = py.ParseV(args, &price, &amount, &msgs)
	if err != nil {
		return py.IncNone(), err
	}
	ID, err := e.api.Sell(price, amount, msgs)
	if err != nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(ID)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetOrder ...
func (e *ExchangePythonLink) GetOrder(args *py.Tuple) (ret *py.Base, err error) {
	var ID string
	err = py.ParseV(args, &ID)
	if err != nil {
		return py.IncNone(), err
	}
	order, err := e.api.GetOrder(ID)
	if err != nil {
		return py.IncNone(), errors.New("get order fail")
	}
	val, ok := pyutil.NewVar(*order)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetOrders ...
func (e *ExchangePythonLink) GetOrders(args *py.Tuple) (ret *py.Base, err error) {
	orders, err := e.api.GetOrders()
	if err != nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(orders)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// CancelOrder ...
func (e *ExchangePythonLink) CancelOrder(args *py.Tuple) (ret *py.Base, err error) {
	var ID string
	err = py.ParseV(args, &ID)
	if err != nil {
		return py.IncNone(), errors.New("cancel order fail")
	}
	_, err = e.api.CancelOrder(ID)
	if err != nil {
		return py.IncNone(), err
	}
	return py.IncNone(), nil
}

// GetTicker ...
func (e *ExchangePythonLink) GetTicker(args *py.Tuple) (ret *py.Base, err error) {
	ticker, err := e.api.GetTicker()
	if err != nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(*ticker)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetRecords ...
func (e *ExchangePythonLink) GetRecords(args *py.Tuple) (ret *py.Base, err error) {
	var period string
	var ma string
	err = py.ParseV(args, &period, &ma)
	if err != nil {
		return py.IncNone(), err
	}
	records, err := e.api.GetRecords(period, ma)
	if err != nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(&records)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// SetContractType ...
func (e *ExchangePythonLink) SetContractType(args *py.Tuple) (ret *py.Base, err error) {
	var s string
	err = py.Parse(args, &s)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetContractType(s)
	return py.IncNone(), nil
}

// GetContractType ...
func (e *ExchangePythonLink) GetContractType(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetContractType()
	val, ok := pyutil.NewVar(s)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// SetDirection ...
func (e *ExchangePythonLink) SetDirection(args *py.Tuple) (ret *py.Base, err error) {
	var s string
	err = py.Parse(args, &s)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetDirection(s)
	return py.IncNone(), nil
}

// GetDirection ...
func (e *ExchangePythonLink) GetDirection(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetDirection()
	val, ok := pyutil.NewVar(s)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// SetMarginLevel ...
func (e *ExchangePythonLink) SetMarginLevel(args *py.Tuple) (ret *py.Base, err error) {
	var s float64
	err = py.Parse(args, &s)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetMarginLevel(s)
	return py.IncNone(), nil
}

// GetMarginLevel ...
func (e *ExchangePythonLink) GetMarginLevel(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetMarginLevel()
	val, ok := pyutil.NewVar(s)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// SetStockType ...
func (e *ExchangePythonLink) SetStockType(args *py.Tuple) (ret *py.Base, err error) {
	var s string
	err = py.Parse(args, &s)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetStockType(s)
	return py.IncNone(), nil
}

// GetStockType ...
func (e *ExchangePythonLink) GetStockType(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetStockType()
	val, ok := pyutil.NewVar(s)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetPosition ...
func (e *ExchangePythonLink) GetPosition(args *py.Tuple) (ret *py.Base, err error) {
	position, err := e.api.GetPosition()
	if err != nil {
		return py.IncNone(), err
	}
	val, ok := pyutil.NewVar(&position)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil

}

// GetBackAccount ...
func (e *ExchangePythonLink) GetBackAccount(args *py.Tuple) (ret *py.Base, err error) {
	account := e.api.GetBackAccount()
	if err != nil {
		return py.IncNone(), errors.New("get backaccount fail")
	}
	val, ok := pyutil.NewVar(account)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// SetBackAccount ...
func (e *ExchangePythonLink) SetBackAccount(args *py.Tuple) (ret *py.Base, err error) {
	var key string
	var val float64
	err = py.ParseV(args, &key, &val)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
	}
	e.api.SetBackAccount(key, val)
	return py.IncNone(), nil
}

// SetBackCommission ...
func (e *ExchangePythonLink) SetBackCommission(args *py.Tuple) (ret *py.Base, err error) {
	var start, end, period float64
	err = py.ParseV(args, &start, &end, &period)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
	}
	e.api.SetBackCommission(start, end, period, period, true)
	return py.IncNone(), nil
}

// GetBackCommission ...
func (e *ExchangePythonLink) GetBackCommission(args *py.Tuple) (ret *py.Base, err error) {
	var commission BackCommission
	taker, maker, rate, cover, coin := e.api.GetBackCommission()
	commission.Taker = taker
	commission.Maker = maker
	commission.ContractRate = rate
	commission.Coin = coin
	commission.CoverRate = cover
	val, ok := pyutil.NewVar(commission)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// SetBackTime ...
func (e *ExchangePythonLink) SetBackTime(args *py.Tuple) (ret *py.Base, err error) {
	var start, end int64
	var period string
	err = py.ParseV(args, &start, &end, &period)
	if err != nil {
		return py.IncNone(), err
	}
	e.api.SetBackTime(start, end, period)
	return py.IncNone(), nil
}

// GetBackTime ...
func (e *ExchangePythonLink) GetBackTime(args *py.Tuple) (ret *py.Base, err error) {
	var backTime BackTime
	start, end, period := e.api.GetBackTime()
	backTime.Start = start
	backTime.End = end
	backTime.Period = period
	val, ok := pyutil.NewVar(backTime)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}
