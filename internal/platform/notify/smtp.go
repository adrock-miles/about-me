package notify

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/adrock-miles/about-me/internal/platform/config"
)

// SMTPSender delivers via net/smtp with optional PLAIN auth (port 587 STARTTLS).
type SMTPSender struct {
	host     string
	port     int
	username string
	password string
}

// NewSMTPSender builds an SMTPSender from config.
func NewSMTPSender(cfg config.SMTPConfig) *SMTPSender {
	return &SMTPSender{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
	}
}

// SendEmail serializes a single text/plain message and hands it to net/smtp.
func (s *SMTPSender) SendEmail(_ context.Context, from, to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body,
	)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
