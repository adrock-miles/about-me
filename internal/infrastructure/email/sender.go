package email

import (
	"fmt"
	"net/smtp"

	"github.com/adrock-miles/about-me/internal/domain/contact"
)

// Sender sends emails via SMTP.
type Sender struct {
	host string
	port int
	from string
	to   string
}

// NewSender creates a new email sender.
func NewSender(host string, port int, from, to string) *Sender {
	return &Sender{host: host, port: port, from: from, to: to}
}

// Send delivers a contact message via email.
func (s *Sender) Send(msg *contact.Message) error {
	subject := msg.Subject
	if subject == "" {
		subject = fmt.Sprintf("Portfolio Contact from %s", msg.Name)
	}

	body := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\n\nReply to: %s",
		msg.Name, s.from, s.to, subject, msg.Body, msg.Email,
	)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, nil, s.from, []string{s.to}, []byte(body))
}
