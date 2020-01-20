package main

import (
	"fmt"
	"github.com/xiyanxiyan10/stockdb/sdk"
	"github.com/xiyanxiyan10/stockdb/types"
)

const (
	uri    = "http://localhost:8765"
	auth   = "username:password"
	market = "haobtc"
	symbol = "BTC/CNY"
)

func main() {
	cli := sdk.NewClient(uri, auth)
	fmt.Printf("%+v\n", cli.GetStats())
	resp := cli.GetMarkets()
	for _, market := range resp.Data {
		symbols := cli.GetSymbols(market).Data
		fmt.Printf("Symbols of %s: %+v\n", market, symbols)
		for _, symbol := range symbols {
			fmt.Printf("MinPeriod of %s: %+v\n", symbol, cli.GetPeriodRange(types.Option{
				Market: market,
				Symbol: symbol,
			}).Data)
		}
	}
}
