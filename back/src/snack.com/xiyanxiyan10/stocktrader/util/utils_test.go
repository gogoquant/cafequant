package util

import (
	"testing"
	"time"
)

func TestTimeConvert(t *testing.T) {
	now := time.Now()
	timeUnix := now.Unix()
	timestr := now.Format(TimeLayout)
	t.Logf("time in:%s %d\n", timestr, timeUnix)
	out, err := TimeStr2Unix(timestr)
	if err != nil {
		t.Fatalf("convert timestr to unix:%s\n", err.Error())
	}
	t.Logf("time out:%d\n", out)

}
