package main

import (
	"fmt"
	"os"

	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/handler"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("command args invalid")
		return
	}
	if err := config.Init(os.Args[1]); err != nil {
		fmt.Printf("config init error is %s\n", err.Error())
		return
	}
	if err := model.Init(); err != nil {
		fmt.Printf("model init error is %s\n", err.Error())
		return
	}
	handler.Init()
}
