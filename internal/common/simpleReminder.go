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

// SendSimpleReminder sends a simple reminder to users from .env (for Dane's 10am reminder)
func SendSimpleReminder(sess *discordgo.Session) {
	fmt.Println("Simple reminder service started.")

	if sess == nil {
		fmt.Println("Error: Discord session is nil")
		return
	}

	// Get reminder message from env
	reminderMsg := os.Getenv("reminder_msg")
	if reminderMsg == "" {
		reminderMsg = "Hello pretti! It's time to take your meds❤️❤️❤️"
	}

	// Extract user IDs from env variable
	userIDs := strings.Split(os.Getenv("user"), ",")

	for _, userID := range userIDs {
		userID = strings.TrimSpace(userID)
		if userID == "" {
			continue
		}

		// Send Discord DM asynchronously
		go func(uid string) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Discord DM panic: %v\n", r)
				}
			}()

			channel, err := sess.UserChannelCreate(uid)
			if err != nil {
				fmt.Printf("Error creating DM channel with user %s: %v\n", uid, err)
				return
			}

			_, err = sess.ChannelMessageSend(channel.ID, reminderMsg)
			if err != nil {
				fmt.Printf("Error sending message to user %s: %v\n", uid, err)
			} else {
				fmt.Printf("Sent simple reminder to user %s\n", uid)
			}
		}(userID)
	}

	// Send email asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Email sending panic: %v\n", r)
			}
		}()

		emailsStr := os.Getenv("EMAIL_TO")
		if emailsStr == "" {
			return
		}

		emailList := strings.Split(emailsStr, ",")
		var trimmedEmails []string
		for _, email := range emailList {
			if trimmed := strings.TrimSpace(email); trimmed != "" {
				trimmedEmails = append(trimmedEmails, trimmed)
			}
		}

		if len(trimmedEmails) == 0 {
			return
		}

		subj := os.Getenv("EMAIL_SUBJECT")
		if subj == "" {
			subj = "Medication reminder"
		}

		if err := sendSimpleEmail(trimmedEmails, subj, reminderMsg); err != nil {
			fmt.Printf("Error sending email to %v: %v\n", trimmedEmails, err)
		} else {
			fmt.Printf("Sent simple email reminder to %v\n", trimmedEmails)
		}
	}()
}

// sendSimpleEmail helper for simple reminders
func sendSimpleEmail(to []string, subject, body string) error {
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

	if port == "465" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		}
		return e.SendWithTLS(addr, auth, tlsConfig)
	}

	return e.Send(addr, auth)
}
