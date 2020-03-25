package pythonbind

import (
	"testing"
)

func TestSum(t *testing.T) {
	go2python := Go2Python{fileName: "/tmp/test.py"}
	if err := go2python.Run(); err != nil {
		t.Errorf("run python error %s\n", err.Error())
		return
	}
	return
}
