package mailer

import "log/slog"

type Mailer interface {
	SendPasswordReset(toEmail, username, resetURL string) error
}

type LogMailer struct {
	logger *slog.Logger
}

func NewLogMailer(logger *slog.Logger) *LogMailer {
	return &LogMailer{logger: logger}
}

func (m *LogMailer) SendPasswordReset(toEmail, username, resetURL string) error {
	m.logger.Info("password reset email (console — no SMTP configured)",
		slog.String("to", toEmail),
		slog.String("username", username),
		slog.String("reset_url", resetURL),
	)
	return nil
}
