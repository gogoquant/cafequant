package gobacktest

import (
	"time"
)

// setup mock for an event
type testEvent struct {
}

func (t testEvent) Time() time.Time {
	return time.Now()
}

func (t *testEvent) SetTime(time time.Time) {
}

func (t testEvent) Symbol() string {
	return "testEvent"
}

func (t *testEvent) SetSymbol(s string) {
}
