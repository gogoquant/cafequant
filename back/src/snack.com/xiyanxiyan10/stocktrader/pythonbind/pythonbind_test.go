package pythonbind

import (
	"testing"
)

func TestSum(t *testing.T) {
	go2python := new(Go2Python)
	go2python.fileName = "/tmp/test.py"
	if err := go2python.Run(); err != nil {
		t.Errorf("run python error %s\n", err.Error())
		return
	}
	return
}
