package email

import (
	"fmt"
	"net/smtp"
)

type Sender struct {
	host string
	port string
	user string
	pass string
	from string
}

func New(host, port, user, pass, from string) *Sender {
	return &Sender{host: host, port: port, user: user, pass: pass, from: from}
}

func (s *Sender) IsConfigured() bool {
	return s.user != "" && s.pass != "" && s.from != ""
}

func (s *Sender) SendWelcome(to string) error {
	if !s.IsConfigured() {
		return nil
	}

	auth := smtp.PlainAuth("", s.user, s.pass, s.host)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"From: Smart Parking <%s>\r\n"+
			"Subject: Welcome to Smart Parking!\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"Hi %s,\r\n\r\n"+
			"Your Smart Parking account has been created successfully.\r\n\r\n"+
			"You can now log in and add your vehicles.\r\n\r\n"+
			"Smart Parking Team\r\n",
		to, s.from, to,
	))

	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{to}, msg)
}
