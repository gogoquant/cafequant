package notice

import (
	"crypto/tls"

	"github.com/go-gomail/gomail"
	log "gopkg.in/logger.v1"
)

// MailHandler ...
type MailHandler interface {
	Set(to string, server string, port int, username string, password string)
	Send(msg string) error
}

// sendMail send change password mail
func sendMail(msg, target, Server string, Port int, Username, Password string, SSL bool) error {
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

// MailServer ...
type MailServer struct {
	server   string
	port     int
	to       string
	username string
	password string
	SSL      bool
	interval int
}

// NewMailHandler ...
func NewMailHandler() MailHandler {
	return &MailServer{
		SSL: false,
	}
}

// Set ...
func (s *MailServer) Set(to, server string, port int, username string, password string) {
	s.server = server
	s.port = port
	s.to = to
	s.username = username
	s.password = password
}

// Send ...
func (s *MailServer) Send(msg string) error {
	return sendMail(msg, s.to, s.server, s.port, s.username, s.password, s.SSL)
}
