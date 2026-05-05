package notify

import (
	"context"
	"log/slog"
)

// logSender prints messages via slog instead of sending them. Used in dev
// when no provider is configured so the contact form remains exercisable
// without real credentials.
type logSender struct {
	logger *slog.Logger
}

func newLogSender(logger *slog.Logger) *logSender {
	return &logSender{logger: logger}
}

func (l *logSender) SendEmail(_ context.Context, from, to, subject, body string) error {
	l.logger.Info("email (dev mode, not actually sent)",
		"from", from,
		"to", to,
		"subject", subject,
		"body", body,
	)
	return nil
}
