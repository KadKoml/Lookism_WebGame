package service

import (
	"fmt"
	"net/smtp"
	"os"
)

type SMTPSender interface {
	SendEmail(to, subject, body string) error
}

type smtpSender struct {
	host string
	port string
}

func NewSMTPSender(host, port string) SMTPSender {
	return &smtpSender{
		host: host,
		port: port,
	}
}

func (s *smtpSender) SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")

	// Fallbacks for testing
	if from == "" {
		from = "example@gmail.com" // User will need to set this in their env
	}
	if password == "" {
		password = "pegd whms sgtx qcmy" // User's provided App Password
	}

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, from, subject, body))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	
	var auth smtp.Auth
	// If it's mailpit (localhost), auth might be nil. For Gmail, use PlainAuth.
	if s.host == "smtp.gmail.com" {
		auth = smtp.PlainAuth("", from, password, s.host)
	}

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
