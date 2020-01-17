package notice

import (
	"fmt"
	"testing"
	"time"
)

func TestMailServer(t *testing.T) {
	server := NewMailServer(10, 5)
	go func() {
		var i int = 0
		for {
			i++
			err := server.Send("hello:"+fmt.Sprintf("%d", i), "873706510@qq.com")
			if err != nil {
				fmt.Printf("msg send error:%s\n", err.Error())
			}
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()

	go func() {
		for {
			err := server.Mail()
			if err != nil {
				fmt.Printf("mail send error:%s\n", err.Error())
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	for {
		time.Sleep(time.Duration(60) * time.Minute)
	}
}
