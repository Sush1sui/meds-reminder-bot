package common

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jordan-wright/email"
)

// RemindUserManual sends all active medications regardless of time (for manual testing)
func RemindUserManual(sess *discordgo.Session) {
	fmt.Println("Manual reminder service started.")

	// Nil check for session
	if sess == nil {
		fmt.Println("Error: Discord session is nil")
		return
	}

	// Load medication schedule
	schedule, err := LoadMedicationState()
	if err != nil {
		fmt.Printf("Error loading medication state: %v\n", err)
		return
	}

	// Update medication counts based on elapsed days
	UpdateMedicationCounts(schedule)

	// Get ALL active medications for manual testing
	var reminders []Medication
	for _, med := range schedule.Medications {
		if med.Active {
			reminders = append(reminders, med)
		}
	}

	if len(reminders) == 0 {
		fmt.Println("No active medications found.")
		return
	}

	// Format the reminder message
	reminderMsg := FormatReminderMessage(reminders)

	// Send to JP (hardcoded user for medication tracking)
	userID := JPDiscordID
	if userID == "" {
		fmt.Println("Error: User ID is empty")
		return
	}

	// Send Discord DM asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Discord DM panic: %v\n", r)
			}
		}()

		channel, err := sess.UserChannelCreate(userID)
		if err != nil {
			fmt.Printf("Error creating DM channel with user %s: %v\n", userID, err)
			return
		}

		_, err = sess.ChannelMessageSend(channel.ID, reminderMsg)
		if err != nil {
			fmt.Printf("Error sending message to user %s: %v\n", userID, err)
		} else {
			fmt.Printf("Sent manual reminder to user %s\n", userID)
		}
	}()

	// Send emails asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Email sending panic: %v\n", r)
			}
		}()

		// Parse and trim email addresses
		emailList := strings.Split(JPEmails, ",")
		var trimmedEmails []string
		for _, email := range emailList {
			if trimmed := strings.TrimSpace(email); trimmed != "" {
				trimmedEmails = append(trimmedEmails, trimmed)
			}
		}

		if len(trimmedEmails) == 0 {
			fmt.Println("No valid email addresses")
			return
		}

		subj := "Medication Reminder (Manual Test)"
		plainText := strings.ReplaceAll(reminderMsg, "**", "")

		if err := sendEmail(trimmedEmails, subj, plainText); err != nil {
			fmt.Printf("Error sending email to %v: %v\n", trimmedEmails, err)
		} else {
			fmt.Printf("Sent manual email reminder to %v\n", trimmedEmails)
		}
	}()
}

func RemindUser(sess *discordgo.Session) {
	fmt.Println("Reminder service started.")

	// Nil check for session
	if sess == nil {
		fmt.Println("Error: Discord session is nil")
		return
	}

	// Load medication schedule
	schedule, err := LoadMedicationState()
	if err != nil {
		fmt.Printf("Error loading medication state: %v\n", err)
		return
	}

	// Update medication counts based on elapsed days
	UpdateMedicationCounts(schedule)

	// Get current reminders for this time
	loc := time.FixedZone("UTC+8", 8*3600)
	now := time.Now().In(loc)
	reminders := GetCurrentReminders(schedule, now)

	if len(reminders) == 0 {
		fmt.Println("No medications due at this time.")
		return
	}

	// Format the reminder message
	reminderMsg := FormatReminderMessage(reminders)

	// Send to JP (hardcoded user for medication tracking)
	userID := JPDiscordID
	if userID == "" {
		fmt.Println("Error: User ID is empty")
		return
	}

	// Send Discord DM asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Discord DM panic: %v\n", r)
			}
		}()

		channel, err := sess.UserChannelCreate(userID)
		if err != nil {
			fmt.Printf("Error creating DM channel with user %s: %v\n", userID, err)
			return
		}

		_, err = sess.ChannelMessageSend(channel.ID, reminderMsg)
		if err != nil {
			fmt.Printf("Error sending message to user %s: %v\n", userID, err)
		} else {
			fmt.Printf("Sent reminder to user %s\n", userID)
		}
	}()

	// Send emails asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Email sending panic: %v\n", r)
			}
		}()

		// Parse and trim email addresses
		emailList := strings.Split(JPEmails, ",")
		var trimmedEmails []string
		for _, email := range emailList {
			if trimmed := strings.TrimSpace(email); trimmed != "" {
				trimmedEmails = append(trimmedEmails, trimmed)
			}
		}

		if len(trimmedEmails) == 0 {
			fmt.Println("No valid email addresses")
			return
		}

		subj := "Medication Reminder"
		plainText := strings.ReplaceAll(reminderMsg, "**", "")

		if err := sendEmail(trimmedEmails, subj, plainText); err != nil {
			fmt.Printf("Error sending email to %v: %v\n", trimmedEmails, err)
		} else {
			fmt.Printf("Sent email reminder to %v\n", trimmedEmails)
		}
	}()
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