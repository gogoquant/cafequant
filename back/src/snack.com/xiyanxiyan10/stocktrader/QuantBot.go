package main

import (
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/handler"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

func main() {
	config.Init()
	model.Init()
	handler.Init()
}
