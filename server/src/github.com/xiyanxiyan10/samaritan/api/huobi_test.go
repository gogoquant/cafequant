package api

import (
	"github.com/huobiapi/REST-GO-demo/services"
	"testing"
	"fmt"
)

func TestHuobi(t *testing.T) {
	api := services.NewHuobiApi("", "")
	ticker, err := api.GetTicker("btcoin")
	if nil != err{
		fmt.Errorf("get ticker fail (%s)", err.Error())
		return
	}
	fmt.Printf("get ticker success (%v)", ticker)
}
