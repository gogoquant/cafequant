package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stockdb/constant"
	"snack.com/xiyanxiyan10/stockdb/sdk"
	"snack.com/xiyanxiyan10/stockdb/types"
)

const (
	uri    = "http://localhost:8765"
	auth   = "username:password"
	market = "haobtc"
	symbol = "BTC/CNY"
)

func main() {
	cli := sdk.NewClient(uri, auth)
	opt := types.Option{Period: constant.Hour, Symbol: symbol, Market: market}
	fmt.Printf("%+v\n", cli.GetTimeRange(opt))
	fmt.Printf("%+v\n", cli.GetOHLCs(opt))
}