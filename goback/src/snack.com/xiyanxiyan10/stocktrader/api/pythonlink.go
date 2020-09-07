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
	Period int64
}

// BackCommission ...
type BackCommission struct {
	Taker        float64
	Maker        float64
	ContractRate float64 // 合约每张价值
}

// ExchangePythonLink ...
type ExchangePythonLink struct {
	api Exchange
}

// Ready ...
func (e *ExchangePythonLink) Ready(args *py.Tuple) (ret *py.Base, err error) {
	e.api.Ready(nil)
	return py.IncNone(), nil
}

// SetIO ...
func (e *ExchangePythonLink) SetIO(args *py.Tuple) (ret *py.Base, err error) {
	var i int
	err = py.ParseV(args, &i)
	if err != nil {
		return
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
		return
	}
	return py.IncNone(), nil
}

// Log ...
func (e *ExchangePythonLink) Log(args *py.Tuple) (ret *py.Base, err error) {
	var msgs []interface{}
	err = py.ParseV(args, &msgs)
	if err != nil {
		return
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
		return
	}
	e.api.Log(vars...)
	return py.IncNone(), nil
}

// Sleep ...
func (e *ExchangePythonLink) Sleep(args *py.Tuple) (ret *py.Base, err error) {
	var action int64
	err = py.ParseV(args, &action)
	if err != nil {
		return
	}
	e.api.Sleep(action)
	return py.IncNone(), nil
}

// GetAccount ...
func (e *ExchangePythonLink) GetAccount(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetAccount()
	if s == nil {
		return py.IncNone(), errors.New("get account fail")
	}
	account := s.(constant.Account)
	val, ok := pyutil.NewVar(account)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetDepth ...
func (e *ExchangePythonLink) GetDepth(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetDepth()
	if s == nil {
		return py.IncNone(), errors.New("get depth fail")
	}
	depth := s.(constant.Depth)
	val, ok := pyutil.NewVar(depth)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// Buy ...
func (e *ExchangePythonLink) Buy(args *py.Tuple) (ret *py.Base, err error) {
	var price, amount string
	var msgs []string
	err = py.ParseV(args, &price, &amount, &msgs)
	if err != nil {
		return
	}
	s := e.api.Buy(price, amount, msgs)
	if s == nil {
		return py.IncNone(), errors.New("buy fail")
	}
	ID := s.(string)
	val, ok := pyutil.NewVar(ID)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// Sell ...
func (e *ExchangePythonLink) Sell(args *py.Tuple) (ret *py.Base, err error) {
	var price, amount string
	var msgs []string
	err = py.ParseV(args, &price, &amount, &msgs)
	if err != nil {
		return
	}
	s := e.api.Sell(price, amount, msgs)
	if s == nil {
		return py.IncNone(), errors.New("sell fail")
	}
	ID := s.(string)
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
		return
	}
	s := e.api.GetOrder(ID)
	if s == nil {
		return py.IncNone(), errors.New("get order fail")
	}
	order := s.(constant.Order)
	val, ok := pyutil.NewVar(order)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetOrders ...
func (e *ExchangePythonLink) GetOrders(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetOrders()
	if s == nil {
		return py.IncNone(), errors.New("get orders fail")
	}
	orders := s.([]constant.Order)
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
	s := e.api.CancelOrder(ID)
	if s == nil {
		return py.IncNone(), errors.New("cancel order fail")
	}
	return py.IncNone(), nil
}

// GetTicker ...
func (e *ExchangePythonLink) GetTicker(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetTicker()
	if s == nil {
		return py.IncNone(), errors.New("get ticker fail")
	}
	ticker := s.(constant.Ticker)
	val, ok := pyutil.NewVar(ticker)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// GetRecords ...
func (e *ExchangePythonLink) GetRecords(args *py.Tuple) (ret *py.Base, err error) {
	var period string
	err = py.ParseV(args, &period)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
	}
	s := e.api.GetRecords(period)
	if s == nil {
		return py.IncNone(), errors.New("get records fail")
	}
	records := s.([]constant.Record)
	val, ok := pyutil.NewVar(records)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil

}

// SetContractType ...
func (e *ExchangePythonLink) SetContractType(args *py.Tuple) (ret *py.Base, err error) {
	var s string
	err = py.ParseV(args, &s)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
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
	err = py.ParseV(args, &s)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
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
	err = py.ParseV(args, &s)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
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
	err = py.ParseV(args, &s)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
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
	s := e.api.GetPosition()
	if s == nil {
		return py.IncNone(), errors.New("get position fail")
	}
	records := s.([]constant.Position)
	val, ok := pyutil.NewVar(records)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil

}

func (e *ExchangePythonLink) GetBackAccount(args *py.Tuple) (ret *py.Base, err error) {
	s := e.api.GetBackAccount()
	if s == nil {
		return py.IncNone(), errors.New("get backaccount fail")
	}
	account := s.(map[string]float64)
	val, ok := pyutil.NewVar(records)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

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

func (e *ExchangePythonLink) SetBackCommission(args *py.Tuple) (ret *py.Base, err error) {
	var start, end, period float64
	err = py.ParseV(args, &start, &end, &period)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
	}
	e.api.SetBackCommission(start, end, period)
	return py.IncNone(), nil
}

func (e *ExchangePythonLink) GetBackCommission(args *py.Tuple) (ret *py.Base, err error) {
	var commission BackCommission
	taker, maker, rate := e.api.GetBackCommission()
	commission.Taker = taker
	commission.Maker = maker
	commission.ContractRate = rate
	val, ok := pyutil.NewVar(commission)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

func (e *ExchangePythonLink) SetBackTime(args *py.Tuple) (ret *py.Base, err error) {
	var start, end, period int64
	err = py.ParseV(args, &start, &end, &period)
	if err != nil {
		return py.IncNone(), errors.New("parse param fail")
	}
	e.api.SetBackTime(start, end, period)
	return py.IncNone(), nil
}

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
