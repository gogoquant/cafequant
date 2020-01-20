package main

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/xiyanxiyan10/stockdb/sdk"
	"github.com/xiyanxiyan10/stockdb/types"
)

const (
	uri    = "http://localhost:8765"
	auth   = "username:password"
	market = "haobtc"
	symbol = "BTC/CNY"
)

var location *time.Location

func main() {
	if loc, err := time.LoadLocation("Asia/Shanghai"); err != nil || loc == nil {
		location = time.Local
	} else {
		location = loc
	}
	opt := types.Option{
		Market: market,
		Symbol: symbol,
	}
	for {
		fetch(opt)
		time.Sleep(10 * time.Minute)
	}
}

type Recrod struct {
	Volume float64 `json:"amount"`
	Price  float64 `json:"price"`
	Side   string  `json:"sid"`
	Time   int64   `json:"ts"`
}

type Records struct {
	Status  int      `json:"status"`
	Records []Recrod `json:"data"`
}

func fetch(opt types.Option) {
	var data Records
	req := httplib.Get("https://api.fmex.com/v2/market/trades/btcusd_p")
	if err := req.ToJSON(&data); err != nil {
		log.Println("parse json error: ", err)
	} else {
		orders := []types.Order{}
		for _, record := range data.Records {
			orders = append(orders, types.Order{
				ID:     fmt.Sprint(record.Volume, "@", record.Price),
				Time:   record.Time,
				Price:  record.Price,
				Amount: record.Volume,
				Type:   record.Side,
			})
		}
		if len(orders) > 0 {
			cli := sdk.NewClient(uri, auth)
			if resp := cli.PutOrders(orders, opt); !resp.Success {
				log.Println("PutOrders error: ", resp.Message)
			} else {
				log.Println("PutOrders successfully")
			}
		}
	}
}
