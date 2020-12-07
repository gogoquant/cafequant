package sdkctp

import (
	"testing"
	"time"
)

func TestCtp(t *testing.T) {
	ctp := NewCtp()
	mdFront := []string{"tcp://180.168.146.187:10110", "tcp://180.168.146.187:10111", "tcp://218.202.237.33:10112"}
	traderFront := []string{"tcp://180.168.146.187:10110", "tcp://180.168.146.187:10111", "tcp://218.202.237.33:10112"}
	ctp.SetTradeAccount(mdFront, traderFront, "9999", "125944", "hongwei", "simnow_client_test", "0000000000000000", "/tmp/stream")
	err := ctp.Start()
	if err != nil {
		t.Log("ctp start err:" + err.Error())
		return
	}

	time.Sleep(30 * time.Second)
	for {
		ctp.ReqQryOrder()
		time.Sleep(10 * time.Second)
	}

	//ctp.SubscribeMarketData([]string{"cu2109", "al2105"})
	//for {
	//	time.Sleep(30 * time.Second)
	//t.Log(ctp.GetFuturesList())
	//}
}
