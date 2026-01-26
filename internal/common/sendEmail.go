package common

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
)

func SendEmail(to string, message string) error {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	user := strings.TrimSpace(os.Getenv("SMTP_USER"))
	pass := strings.TrimSpace(os.Getenv("SMTP_PASS"))
	from := strings.TrimSpace(os.Getenv("EMAIL_FROM"))

	if from == "" {
        from = user
    }
    if host == "" || port == "" || user == "" || pass == "" {
        return fmt.Errorf("SMTP config missing")
    }

    addr := host + ":" + port
    e := email.NewEmail()
    e.From = from
    e.To = []string{to}
    e.Subject = "Reminder from Sushi Reminder"
    e.Text = []byte(message)
    e.HTML = []byte(message)

    auth := smtp.PlainAuth("", user, pass, host)

    // Use implicit TLS for port 465, otherwise use STARTTLS (e.Send)
    if port == "465" {
        tlsConfig := &tls.Config{
            InsecureSkipVerify: false,
            ServerName:         host,
        }
        return e.SendWithTLS(addr, auth, tlsConfig)
    }

    // For 587 (STARTTLS) and other non-implicit-TLS ports
    return e.Send(addr, auth)
}