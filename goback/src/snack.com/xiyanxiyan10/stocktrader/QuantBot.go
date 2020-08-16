package main

import (
	"snack.com/xiyanxiyan10/stocktrader/handler"
	"snack.com/xiyanxiyan10/stocktrader/plugin"
)

func main() {
	plugin.Load()
	handler.Server()
}
