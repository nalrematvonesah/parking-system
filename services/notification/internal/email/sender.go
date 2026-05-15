package email

import (
	"fmt"
	"net/smtp"

	"go.uber.org/zap"
)

type Sender struct {
	host string
	port string
	user string
	pass string
	from string
	log  *zap.Logger
}

func New(host, port, user, pass, from string, log *zap.Logger) *Sender {
	return &Sender{host: host, port: port, user: user, pass: pass, from: from, log: log}
}

func (s *Sender) configured() bool {
	return s.user != "" && s.pass != "" && s.from != ""
}

// Send delivers an email. If SMTP isn't configured, the email is logged
// instead of sent — handy for demos without real SMTP credentials.
func (s *Sender) Send(to, subject, body string) error {
	if !s.configured() {
		s.log.Info("EMAIL (dev mode, SMTP not configured)",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.String("body", body),
		)
		return nil
	}

	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"From: Smart Parking <%s>\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		to, s.from, subject, body,
	))

	if err := smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{to}, msg); err != nil {
		s.log.Warn("smtp send failed", zap.String("to", to), zap.Error(err))
		return err
	}
	s.log.Info("email sent", zap.String("to", to), zap.String("subject", subject))
	return nil
}
