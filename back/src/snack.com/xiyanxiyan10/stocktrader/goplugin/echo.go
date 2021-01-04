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

// Init ...
func (e *EchoStragey) Init(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Init")
	return nil
}

// Run ...
func (e *EchoStragey) Run(v ...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Run")
	return nil
}

// Exit ...
func (e *EchoStragey) Exit(...interface{}) interface{} {
	e.Logger.Log(constant.INFO, "", 0.0, 0.0, "Exit")
	return nil
}
