package notice

import (
	"github.com/blinkbean/dingtalk"
)

// DingHandler ...
type DingHandler interface {
	Set(token string, key string)
	Send(msg string) error
}

// DingServer ...
type DingServer struct {
	token []string
	key   string
}

// NewDingHandler ...
func NewDingHandler() DingHandler {
	return &DingServer{}
}

// Set ...
func (s *DingServer) Set(token, key string) {
	if token == "" {
		s.token = make([]string, 1)
	} else {
		s.token = append(s.token, token)
	}
	s.key = key
}

// Send ...
func (s *DingServer) Send(msg string) error {
	cli := dingtalk.InitDingTalk(s.token, s.key)
	return cli.SendTextMessage(msg)
}
