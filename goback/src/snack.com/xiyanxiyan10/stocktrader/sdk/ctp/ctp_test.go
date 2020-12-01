package sdkctp

import "testing"

func TestCtp(t *testing.T) {
	ctp := NewCtp()
	ctp.SetTradeAccount([]string{}, []string{}, "", "", "", "", "", "/tmp/stream")
	err := ctp.Start()
	if err != nil {
		t.Log("ctp start err:" + err.Error())
		return
	}
}
