package goplugin

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// EchoStragey ...
type EchoStragey struct {
	GoStragey
}

// NewEchoHandler ...
func NewEchoHandler(...interface{}) (GoStrageyHandler, error) {
	var echo EchoStragey
	return &echo, nil
}

// Set ...
func (e *EchoStragey) Set(string, interface{}) interface{} {
	return nil
}

// Init ...
func (e *EchoStragey) Init(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Init")
	return nil
}

// Call ...
func (e *EchoStragey) Call(name string, v ...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Call")
	return nil
}

// Exit ...
func (e *EchoStragey) Exit(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Exit")
	return nil
}
