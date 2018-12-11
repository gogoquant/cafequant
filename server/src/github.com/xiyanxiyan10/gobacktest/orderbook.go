package gobacktest

import (
	"fmt"
	"sort"
	"sync"
)

func NewOrderBook() OrderBook {
	return OrderBook{
		counter:    0,
		orders:     []OrderEvent{},
		history:    []OrderEvent{},
	}
}

// OrderBook represents an order book.
type OrderBook struct {
	lock       sync.Mutex
	counter    int
	orders     []OrderEvent
	history    []OrderEvent
}

// Add an order to the order book.
func (ob *OrderBook) Add(order OrderEvent) error {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	// increment counter
	ob.counter++
	// assign an ID to the Order
	order.SetID(ob.counter)

	ob.orders = append(ob.orders, order)

	//ob.EnableSubscribe(order.Symbol())
	return nil
}

// Remove an order from the order book, append it to history.
func (ob *OrderBook) Remove(id int) error {
	for i, order := range ob.orders {
		// order found
		if order.ID() == id {

			//ob.DisableSubscribe(order.Symbol())

			ob.history = append(ob.history, ob.orders[i])

			ob.orders = append(ob.orders[:i], ob.orders[i+1:]...)

			return nil
		}
	}

	// order not found
	return fmt.Errorf("order with id %v not found", id)
}

// CancelOrder Remove an order from the order book, append it to history.
func (ob *OrderBook) CancelOrder(id int) error {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	return ob.Remove(id)
}

// CommitOrder ...
func (ob *OrderBook) CommitOrder(id int) (*Fill, error) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	for _, order := range ob.orders {
		// order found
		if order.ID() == id {
			if order.Status() == OrderCanceled || order.Status() == OrderCancelPending {
				return nil, fmt.Errorf("order with id %v canceled again", id)
			}
			order.Submit()

			fill := new(Fill)
			fill.SetSymbol(order.Symbol())
			fill.SetTime(order.Time())
			fill.SetQuantifier(order.Quantifier())

			return fill, nil
		}
	}
	return nil, fmt.Errorf("order with id %v not found", id)
}

// Orders returns all Orders from the order book
func (ob *OrderBook) Orders() ([]OrderEvent, bool) {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	orders := ob.deepCopyOrders(ob.orders)
	if len(orders) == 0 {
		return orders, false
	}
	return orders, true
}

// Orders returns all Orders from the order book
func (ob *OrderBook) deepCopyOrders(backorders []OrderEvent) []OrderEvent {
	var orders []OrderEvent
	for _, backorder := range backorders {
		var o Order
		var order OrderEvent
		order = &o
		order.Fill(backorder)
		orders = append(orders, order)
	}
	return orders
}

// OrderBy returns the order by a select function from the order book.
func (ob *OrderBook) OrderBy(fn func(order OrderEvent) bool) ([]OrderEvent, bool) {
	var orders = []OrderEvent{}

	for _, order := range ob.orders {
		if fn(order) {
			orders = append(orders, order)
		}
	}

	if len(orders) == 0 {
		return orders, false
	}

	return orders, true
}

// OrdersBySymbol returns the order of a specific symbol from the order book.
func (ob *OrderBook) OrdersBySymbol(symbol string) ([]OrderEvent, bool) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	var fn = func(order OrderEvent) bool {
		if order.Symbol() != symbol {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return ob.deepCopyOrders(orders), ok
}

// OrdersBidBySymbol returns all bid orders of a specific symbol from the order book.
func (ob *OrderBook) OrdersBidBySymbol(symbol string) ([]OrderEvent, bool) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	var fn = func(order OrderEvent) bool {
		if (order.Symbol() != symbol) || (order.Direction() != BOT) {
			return false
		}
		return true
	}
	orders, ok := ob.OrderBy(fn)

	// sort bid orders ascending, lowest price first
	sort.Slice(orders, func(i, j int) bool {
		o1 := orders[i]
		o2 := orders[j]

		return o1.Limit() < o2.Limit()

	})

	return ob.deepCopyOrders(orders), ok
}

// OrdersAskBySymbol returns all bid orders of a specific symbol from the order book.
func (ob *OrderBook) OrdersAskBySymbol(symbol string) ([]OrderEvent, bool) {
	var fn = func(order OrderEvent) bool {
		if (order.Symbol() != symbol) || (order.Direction() != SLD) {
			return false
		}
		return true
	}
	orders, ok := ob.OrderBy(fn)

	// sort bid orders descending, highest price first
	sort.Slice(orders, func(i, j int) bool {
		o1 := orders[i]
		o2 := orders[j]

		return o1.Limit() > o2.Limit()

	})

	return ob.deepCopyOrders(orders), ok
}

// OrdersOpen returns all orders which are open from the order book.
func (ob *OrderBook) OrdersOpen() ([]OrderEvent, bool) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	var fn = func(order OrderEvent) bool {
		if (order.Status() != OrderNew) || (order.Status() != OrderSubmitted) || (order.Status() != OrderPartiallyFilled) {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return ob.deepCopyOrders(orders), ok
}

// OrdersCanceled returns all orders which are canceled from the order book.
func (ob *OrderBook) OrdersCanceled() ([]OrderEvent, bool) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	var fn = func(order OrderEvent) bool {
		if (order.Status() == OrderCanceled) || (order.Status() == OrderCancelPending) {
			return true
		}
		return false
	}

	orders, ok := ob.OrderBy(fn)
	return ob.deepCopyOrders(orders), ok
}
