package notifications

import (
	"fmt"
	"net/smtp"
	"uptime-monitor/internal/config"
)

// EmailNotifier handles sending emails.
type EmailNotifier struct {
	cfg *config.Config
}

// NewEmailNotifier creates a new notifier.
func NewEmailNotifier(cfg *config.Config) *EmailNotifier {
	return &EmailNotifier{cfg: cfg}
}

// SendNotification sends an email to a recipient.
func (n *EmailNotifier) SendNotification(to, subject, body string) error {
	// Check if SMTP is configured
	if n.cfg.SMTPHost == "" || n.cfg.SMTPPort == "" || n.cfg.SMTPUsername == "" || n.cfg.SMTPPassword == "" {
		return fmt.Errorf("SMTP not configured. Skipping email notification")
	}

	auth := smtp.PlainAuth("", n.cfg.SMTPUsername, n.cfg.SMTPPassword, n.cfg.SMTPHost)

	// Construct the email message
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, n.cfg.EmailSender, subject, body))

	addr := fmt.Sprintf("%s:%s", n.cfg.SMTPHost, n.cfg.SMTPPort)

	return smtp.SendMail(addr, auth, n.cfg.EmailSender, []string{to}, msg)
}
