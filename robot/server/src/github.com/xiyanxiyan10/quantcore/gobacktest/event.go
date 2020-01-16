package gobacktest

import (
	"time"
)

// EventHandler declares the basic event interface
type EventHandler interface {
	Timer
	SymbolHandler
}

// Timer declares the timer interface
type Timer interface {
	Time() time.Time
	SetTime(time.Time)
}

// SymbolHandler declares the symbol interface
type SymbolHandler interface {
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
	IDer
	EventHandler
	SignalPricer
	Quantifier
	DirectionHandler
}

// Directioner defines a direction interface
type DirectionHandler interface {
	Direction() Direction
	SetDirection(Direction)
}

// OrderEvent declares the order event interface.
type OrderEvent interface {
	EventHandler
	DirectionHandler
	Quantifier
	IDer
	SignalPricer
	Status() OrderStatus
	Limit() float64
	Stop() float64
}

// Quantifier defines a qty interface.
type Quantifier interface {
	Qty() float64
	SetQty(float64)
}

// SignalPricer defines a qty interface.
type SignalPricer interface {
	Price() float64
	SetPrice(float64)
}

// IDer declares setting and retrieving of an Id.
type IDer interface {
	ID() int
	SetID(int)
}

// FillEvent declares fill event functionality.
type FillEvent interface {
	EventHandler
	DirectionHandler
	Quantifier
	Price() float64
	Commission() float64
	ExchangeFee() float64
	Cost() float64
	Value() float64
	NetValue() float64
}
