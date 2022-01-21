package main

import (
	"fmt"
	"os"

	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/handler"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

func main() {
	var configPath string
	if len(os.Args) < 2 {
		configPath = "./config.ini"
	} else {
		configPath = os.Args[1]
	}
	if err := config.Init(configPath); err != nil {
		fmt.Printf("config init error is %s\n", err.Error())
		return
	}
	if err := model.Init(); err != nil {
		fmt.Printf("model init error is %s\n", err.Error())
		return
	}
	handler.Init()
}
