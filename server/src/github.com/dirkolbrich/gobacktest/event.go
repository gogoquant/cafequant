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

// OrderEvent declares the order event interface.
type OrderEvent interface {
	EventHandler
	Quantifier
	IDer
	Cancel()
	Submit()
	Status() OrderStatus
	Limit() float64
	Stop() float64
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
	SetQuantifier(orderType OrderType, qtyType QtyType, qty int64, fqty float64, direction Direction, price float64)
	Quantifier() (orderType OrderType, qtyType QtyType, qty int64, fqty float64, direction Direction, price float64)
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
	Cost() float64
	Value() float64
	NetValue() float64
}