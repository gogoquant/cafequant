package gobacktest

import (
	"reflect"
	"testing"
	"time"
)

func TestOnOrder(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	// set ExecutionHandler with symbol
	var e = &Exchange{
		Symbol:      "TEST",
		Commission:  &FixedCommission{Commission: 0},
		ExchangeFee: &FixedExchangeFee{ExchangeFee: 1.0},
	}

	// orderEventTests is a table for testing parsing bar data into a BarEvent
	var testCases = []struct {
		msg     string
		order   OrderEvent  // OrderEvent input
		data    DataHandler // DataEvent input
		expFill FillEvent   // expected FillEvent return
		expErr  error       // expected error output
	}{
		{
			"buy order",
			&Order{
				Event: Event{timestamp: exampleTime, symbol: "TEST.DE"},
				//direction: BOT, // buy or sell
				qty: 10,
			},
			&Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Close: 10},
				},
			},
			&Fill{
				Event:    Event{timestamp: exampleTime, symbol: "TEST.DE"},
				Exchange: "TEST",
				//direction:   BOT, // BOT for buy or SLD for sell
				qty:         10,
				price:       10,
				commission:  0,
				exchangeFee: 1,
				cost:        1,
			},
			nil,
		},
		{
			"sell order",
			&Order{
				Event: Event{timestamp: exampleTime, symbol: "TEST.DE"},
				//direction: SLD, // buy or sell
				qty: 10,
			},
			&Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Close: 10},
				},
			},
			&Fill{
				Event:    Event{timestamp: exampleTime, symbol: "TEST.DE"},
				Exchange: "TEST",
				//direction:   SLD, // BOT for buy or SLD for sell
				qty:         10,
				price:       10,
				commission:  0,
				exchangeFee: 1,
				cost:        1,
			},
			nil,
		},
	}

	for _, tc := range testCases {
		fill, err := e.OnOrder(tc.order, tc.data)
		if !reflect.DeepEqual(fill, tc.expFill) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%s OnOrder(%v): \nexpected %+v %v, \nactual   %+v %v",
				tc.msg, tc.order, tc.expFill, tc.expErr, fill, err)
		}
	}
}
