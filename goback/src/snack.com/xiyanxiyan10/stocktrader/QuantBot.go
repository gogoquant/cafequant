package main

import (
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/handler"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

func main() {
	config.InitConfig()
	model.InitModel()
	handler.Server()
}
