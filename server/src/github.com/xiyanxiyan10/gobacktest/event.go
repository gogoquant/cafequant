package gobacktest

import (
	"time"
)

// EventHandler declares the basic event interface
type EventHandler interface {
	Timer
	Symboler
}

// Timer declares the timer interface
type Timer interface {
	Time() time.Time
	SetTime(time.Time)
}

// Symboler declares the symboler interface
type Symboler interface {
	Symbol() string
	SetSymbol(string)
}

// Event is the implementation of the basic event interface.
type Event struct {
	timestamp time.Time
	symbol    string
}

// Time returns the timestamp of an event
func (e Event) Time() time.Time {
	return e.timestamp
}

// SetTime returns the timestamp of an event
func (e *Event) SetTime(t time.Time) {
	e.timestamp = t
}

// Symbol returns the symbol string of the event
func (e Event) Symbol() string {
	return e.symbol
}

// SetSymbol returns the symbol string of the event
func (e *Event) SetSymbol(s string) {
	e.symbol = s
}

// SignalEvent declares the signal event interface.
type SignalEvent interface {
	EventHandler
	Quantifier
}

// DataGramEvent ...
type DataGramEvent interface {
	EventHandler
	SetId(uid string)
	Id() string
	SetTag(key, val string)
	Tags()map[string]string
	SetColor(c string)
	Color(c string) string
	SetVal(key, v interface{})
	Val() map[string]interface{}
}

// OrderEvent declares the order event interface.
type OrderEvent interface {
	EventHandler
	Quantifier
	IDer
	Cancel()
	Submit()
	Fill(from OrderEvent)
	SetStatus(OrderStatus)
	Status() OrderStatus
	Limit() float64
	Stop() float64
}

type CmdEvent interface {
	EventHandler
	CmdHandler
}

type CmdHandler interface {
	Cmd() string
}

// Quantifier defines a qty interface.
type Quantifier interface {
	Qty() int64
	SetQty(int64)
	FQty() float64
	SetFQty(float64)
	Price() float64
	SetPrice(i float64)
	QtyType() QtyType
	SetQtyType(i QtyType)
	OrderType() OrderType
	SetOrderType(i OrderType)
	Direction() Direction
	SetDirection(dir Direction)
	SetFeeHandler(handler FeeHandler)
	FeeHandler()FeeHandler
	Exchange()(string)
	SetExchange(string)
	SetQuantifier(exchange string, orderType OrderType, qtyType QtyType, qty int64, fqty float64, direction Direction, price float64, feeHandler FeeHandler)
	Quantifier() (exchange string, orderType OrderType, qtyType QtyType, qty int64, fqty float64, direction Direction, price float64, feeHandler FeeHandler)
}

// IDer declares setting and retrieving of an Id.
type IDer interface {
	ID() int
	SetID(int)
}

// FillEvent declares fill event functionality.
type FillEvent interface {
	EventHandler
	Quantifier
	Commission() float64
	ExchangeFee() float64
	SetCost(cost float64)
	Cost() float64
	Value() float64
	NetValue() float64
}
