package mailer

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
)

type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
	logger   *slog.Logger
}

func NewSMTPMailer(host, port, username, password, from string, logger *slog.Logger) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		logger:   logger,
	}
}

func (m *SMTPMailer) SendPasswordReset(toEmail, username, resetURL string) error {
	subject := "GoAdmin — Password Reset"
	body := fmt.Sprintf(`Hello %s,

A password reset was requested for your GoAdmin account.

Click the link below to reset your password:
%s

This link expires in 1 hour.

If you didn't request this, please ignore this email.

— GoAdmin`, username, resetURL)

	return m.send(toEmail, subject, body)
}

func (m *SMTPMailer) send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body)

	var auth smtp.Auth
	if m.username != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}

	var err error
	if auth != nil {
		err = smtp.SendMail(addr, auth, m.from, []string{to}, []byte(msg))
	} else {
		err = smtp.SendMail(addr, nil, m.from, []string{to}, []byte(msg))
	}

	if err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}

	m.logger.Info("password reset email sent",
		slog.String("to", to),
	)

	return nil
}

// New returns the best available mailer.
// If SMTP is configured, returns SMTPMailer; otherwise falls back to LogMailer.
func New(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom string, logger *slog.Logger) Mailer {
	if smtpHost != "" {
		cleanedHost := strings.TrimPrefix(smtpHost, "smtp.")
		return NewSMTPMailer(cleanedHost, smtpPort, smtpUser, smtpPass, smtpFrom, logger)
	}
	return NewLogMailer(logger)
}
