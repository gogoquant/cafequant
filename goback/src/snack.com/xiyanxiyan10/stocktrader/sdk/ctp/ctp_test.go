package sdkctp

import "testing"

func TestCtp(t *testing.T) {
	ctp := NewCtp()
	err := ctp.Start()
	ctp.SetTradeAccount([]string{}, []string{}, "", "", "", "", "", "/tmp/stream")
	if err != nil {
		t.Log("ctp start err:" + err.Error())
		return
	}
}
