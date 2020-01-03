package api

import (
	"encoding/json"
	"fmt"
	"github.com/xiyanxiyan10/quantcore/constant"

	//"github.com/xiyanxiyan10/quantcore/constant"
	"testing"
)

func TestFuture(t *testing.T) {
	var opt constant.Option
	//set key and exchange name
	opt.Type = "fmex"
	opt.Name = "test"
	opt.AccessKey = "e44065942e0c4ae789ace7e1496aa991"
	opt.SecretKey = "5b531a50875c42f094bce1fb94de24e4"
	exchange := NewFutureExchange(opt)
	exchange.Init()
	var exchangeAPI Exchange
	exchangeAPI = exchange

	//Account get
	account := exchangeAPI.GetAccount()
	bytes, err := json.Marshal(account)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	fmt.Printf("Get account:%s\n", string(bytes))

	//type get
	exchangeType := exchangeAPI.GetType()
	fmt.Printf("exchange Type:%s\n", exchangeType)

	//Name get
	exchangeName := exchangeAPI.GetName()
	fmt.Printf("exchange Type:%s\n", exchangeName)

	//Min amount get
	fmt.Printf("min amount %v\n", exchangeAPI.GetMinAmount("BTC/USD"))
	exchangeAPI.SetStockType("BTC/USD")
	fmt.Printf("stockType %v\n", exchangeAPI.GetStockType())

	exchangeAPI.SetContractType("this_week")
	fmt.Printf("contractType %v\n", exchangeAPI.GetContractType())

	ticker := exchangeAPI.GetTicker()
	bytes, err = json.Marshal(ticker)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	fmt.Printf("ticker %v\n", string(bytes))

	//TODO
	records := exchangeAPI.GetRecords("M15", 10, 0)
	bytes, err = json.Marshal(records)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	fmt.Printf("records %v\n", string(bytes))

	// no order Type in fmex
	orders := exchangeAPI.GetOrders()
	bytes, err = json.Marshal(orders)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	fmt.Printf("orders %v\n", string(bytes))

	// no order Type in fmex
	position := exchangeAPI.GetPosition()
	bytes, err = json.Marshal(position)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	fmt.Printf("position %v\n", string(bytes))

	// no order Type in fmex
	depth := exchangeAPI.GetDepth(3)
	bytes, err = json.Marshal(depth)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	fmt.Printf("position %v\n", string(bytes))
	/*
		var orderId interface{}

		exchangeAPI.SetDirection(constant.TradeTypeLong)
		orderId = exchangeAPI.Buy("0", "1")
		bytes, err = json.Marshal(orderId)
		if err != nil {
			fmt.Println("json.Marshal failed:", err)
			return
		}
		fmt.Printf("buy %v\n", string(bytes))

		exchangeAPI.SetDirection(constant.TradeTypeLong)
		orderId = exchangeAPI.Buy("7000", "1")
		bytes, err = json.Marshal(orderId)
		if err != nil {
			fmt.Println("json.Marshal failed:", err)
			return
		}
		fmt.Printf("buy %v\n", string(bytes))


		exchangeAPI.SetDirection(constant.TradeTypeShort)
		orderId = exchangeAPI.Sell("0", "1")
		bytes, err = json.Marshal(orderId)
		if err != nil {
			fmt.Println("json.Marshal failed:", err)
			return
		}
		fmt.Printf("sell %v\n", string(bytes))
	*/
	orderIdStr := fmt.Sprint(10077304541912)
	order := exchangeAPI.GetOrder(orderIdStr)

	bytes, err = json.Marshal(order)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	fmt.Printf("order %v\n", string(bytes))

	cancelSuccess := exchangeAPI.CancelOrder(orderIdStr)
	fmt.Printf("order cancel %v\n", cancelSuccess)

	traders := exchangeAPI.GetTrades()
	fmt.Printf("traders  %v\n", traders)

}
