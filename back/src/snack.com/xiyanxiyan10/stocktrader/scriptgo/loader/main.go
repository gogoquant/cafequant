package main

import "snack.com/xiyanxiyan10/stocktrader/goplugin"

// NewLoaderHandler ...
func NewLoaderHandler(v ...interface{}) (goplugin.GoStrageyHandler, error) {
	return goplugin.NewLoaderHandler(v)
}

func main() {
	goplugin.RunLoader()
}
