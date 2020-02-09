package okcoin

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var okcn = NewOKCoinCn(http.DefaultClient, "", "")

func TestOKCoinCN_API_GetTicker(t *testing.T) {
	ticker, _ := okcn.GetTicker(goex.BTC_CNY)
	t.Log(ticker)
}

func TestOKCoinCN_API_GetDepth(t *testing.T) {
	dep, _ := okcn.GetDepth(1, goex.ETH_CNY)
	t.Log(dep)
}

func TestOKCoinCN_API_GetKlineRecords(t *testing.T) {
	klines, _ := okcn.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN, 1000, -1)
	t.Log(klines)
}
