package notice

import (
	"crypto/tls"
	"errors"
	"fmt"
	log "gopkg.in/logger.v1"
	"time"

	"github.com/go-gomail/gomail"
)

// MailHandler ...
type MailHandler interface {
	Set(server string, port int, username string, password string)
	Send(msg, to string) error
	Start() error
	Stop() error
	Status() string
	Mail() error
}

//SendMail send change password mail
func SendMail(msg, target, Server string, Port int, Username, Password string, SSL bool) error {
	m := gomail.NewMessage(gomail.SetCharset("UTF-8"), gomail.SetEncoding(gomail.Base64))
	m.SetAddressHeader("From", Username, "[量化]") // 发件人
	m.SetHeader("To", target)                    // 收件人

	m.SetHeader("Subject", "量化") // 主题
	m.SetBody("text/html", msg)  // 正文
	log.Debugf("smtp config : host=%s,port=%d", Server, Port)
	d := gomail.NewDialer(Server, Port, Username, Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !SSL}

	if err := d.DialAndSend(m); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// MailMsg ...
type MailMsg struct {
	To  string
	Msg string
}

// MailServer ...
type MailServer struct {
	server    string
	port      int
	run       bool
	username  string
	password  string
	SSL       bool
	cacheSize int
	interval  int
	ch        chan MailMsg
}

// NewMailHandler ...
func NewMailHandler(cacheSize, interval int) MailHandler {
	return &MailServer{
		cacheSize: cacheSize,
		interval:  interval,
		SSL:       false,
		run:       false,
		ch:        make(chan MailMsg, cacheSize),
	}
}

// Set ...
func (s *MailServer) Set(server string, port int, username string, password string) {
	s.server = server
	s.port = port
	s.username = username
	s.password = password
}

// Send ...
func (s *MailServer) Send(msg, to string) error {
	select {
	case s.ch <- MailMsg{Msg: msg, To: to}:
		return nil
	default:
		return errors.New("no message sent")
	}
}

func (s *MailServer) Start() error {
	if s.run == true {
		return errors.New("already start")
	}
	go func() {
		s.run = true
		for {
			if s.run != true {
				break
			}
			err := s.Mail()
			if err != nil {
				fmt.Printf("mail send error:%s\n", err.Error())
			}
			time.Sleep(time.Duration(s.interval) * time.Minute)
		}
	}()
	return nil
}

func (s *MailServer) Stop() error {
	if s.run != true {
		return errors.New("already stop")
	}
	s.run = false
	return nil
}

func (s *MailServer) Status() string {
	if s.run {
		return "run"
	}
	return "stop"
}

func (s *MailServer) Mail() error {
	msg := <-s.ch
	return SendMail(msg.Msg, msg.To, s.server, s.port, s.username, s.password, s.SSL)
}
