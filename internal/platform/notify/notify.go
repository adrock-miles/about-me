// Package notify delivers transactional notifications. For the portfolio it
// only handles contact-form emails, but it's structured so a second channel
// (e.g. Slack) could slot in next to email.
package notify

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/adrock-miles/about-me/internal/platform/config"
)

// Message is a single contact-form submission ready to deliver.
type Message struct {
	Name    string
	Email   string
	Subject string
	Body    string
}

// Render produces the (subject, body) pair that providers send verbatim.
func (m Message) Render() (subject, body string) {
	subject = m.Subject
	if subject == "" {
		subject = fmt.Sprintf("Portfolio contact from %s", m.Name)
	}
	body = fmt.Sprintf("From: %s <%s>\n\n%s\n", m.Name, m.Email, m.Body)
	return subject, body
}

// Sender is the transport contract — anything that can put bytes on the wire.
type Sender interface {
	SendEmail(ctx context.Context, from, to, subject, body string) error
}

// Mailer is the high-level entrypoint handlers use. It wraps a Sender with
// the configured from/to addresses so callers don't carry config around.
type Mailer struct {
	cfg    config.EmailConfig
	sender Sender
	logger *slog.Logger
}

// NewMailer picks a Sender based on cfg.Provider:
//
//	""        → logSender (dev-mode, prints via slog)
//	"smtp"    → SMTPSender
//	"resend"  → ResendSender
//
// An unknown provider falls back to log mode with a warning so a typo can't
// silently swallow contact-form messages in production.
func NewMailer(cfg config.EmailConfig, logger *slog.Logger) *Mailer {
	if logger == nil {
		logger = slog.Default()
	}

	var sender Sender
	switch strings.ToLower(cfg.Provider) {
	case "", "log":
		sender = newLogSender(logger)
	case "smtp":
		sender = NewSMTPSender(cfg.SMTP)
	case "resend":
		sender = NewResendSender(cfg.Resend.APIKey)
	default:
		logger.Warn("unknown email provider; falling back to log",
			"provider", cfg.Provider)
		sender = newLogSender(logger)
	}

	return &Mailer{cfg: cfg, sender: sender, logger: logger}
}

// SendContact delivers a contact-form message to the configured recipient.
func (m *Mailer) SendContact(ctx context.Context, msg Message) error {
	subject, body := msg.Render()
	return m.sender.SendEmail(ctx, m.cfg.From, m.cfg.To, subject, body)
}
