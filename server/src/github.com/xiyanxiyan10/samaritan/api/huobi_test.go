package api

import (
	"fmt"
	"github.com/huobiapi/REST-GO-demo/services"
	"testing"
)

func TestHuobi(t *testing.T) {
	api := services.NewHuobiApi("fa2e1524-c6354b16-a72983d9-965bf", "e220db95-99575247-b157531a-e2e34")
	ticker, err := api.GetTicker("datxbtc")
	if nil != err {
		fmt.Errorf("get ticker fail (%s)", err.Error())
		return
	}
	fmt.Printf("get ticker success (%v)", ticker)
}
