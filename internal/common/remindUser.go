package common

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jordan-wright/email"
)

func RemindUser(sess *discordgo.Session) {
	fmt.Println("Reminder service started.")

	// extract user ids from env variable
	userIDs := strings.Split(os.Getenv("user"), ",")

	for _, userID := range userIDs {
		// create a DM channel with the user
		channel, err := sess.UserChannelCreate(userID)
		if err != nil {
			fmt.Printf("Error creating DM channel with user %s: %v\n", userID, err)
			continue
		}
		// send a reminder message
        reminder := os.Getenv("reminder_msg")
		_, err = sess.ChannelMessageSend(channel.ID, os.Getenv("reminder_msg"))
		if err != nil {
			fmt.Printf("Error sending message to user %s: %v\n", userID, err)
			continue
		}
		fmt.Printf("Sent reminder to user %s\n", userID)

		// send email as well
        emails := strings.Split(os.Getenv("EMAIL_TO"), ",")
        subj := os.Getenv("EMAIL_SUBJECT")
        if subj == "" {
            subj = "Medication reminder"
        }
        if len(emails) > 0 && emails[0] != "" {
            if err := sendEmail(emails, subj, reminder); err != nil {
                fmt.Printf("Error sending email to %v: %v\n", emails, err)
            } else {
                fmt.Printf("Sent email reminder to %v\n", emails)
            }
        }
	}
}


// helper to send email using github.com/jordan-wright/email
func sendEmail(to []string, subject, body string) error {
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
    e.To = to
    e.Subject = subject
    e.Text = []byte(body)
    e.HTML = []byte(body)

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