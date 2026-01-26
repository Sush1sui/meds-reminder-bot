package common

import (
	"fmt"
	"time"

	"github.com/Sush1sui/meds_reminder/internal/repository"
	"github.com/bwmarrin/discordgo"
)

// schedules a repeating daily reminder for the given user at the
// specified hour/minute in AM/PM format
// it spawns a goroutine and returns immediately
// errors are returned only for invalid input or failure to fetch
// initial reminders when calling StartAllReminders.
func StartReminder(s *discordgo.Session, userID, message, email string, hour, minute int, am bool) error {
	if userID == "" || email == "" || message == "" {
		return fmt.Errorf("userID, email, and message cannot be empty")
	}
	if hour < 1 || hour > 12 || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid time: hour must be 1-12, minute 0-59")
	}

	// convert to 24-hour clock
	hour24 := hour % 12
	if !am {
		hour24 += 12
	}

	go func() {
		for {
			now := time.Now()
			loc := now.Location()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour24, minute, 0, 0, loc)
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}

			wait := time.Until(next)
			time.Sleep(wait)

			// send email
			err := SendEmail(email, message)
			if err != nil {
				fmt.Printf("failed to send email to %v: %v\n", email, err)
			}

			// send DM
			ch, err := s.UserChannelCreate(userID)
			if err != nil {
				fmt.Printf("failed to create DM channel for user %v: %v\n", userID, err)
				continue
			}
			_, _ = s.ChannelMessageSend(ch.ID, message)

			// loop will compute next and sleep another 24h (via recompute logic)
		}
	}()

	return nil
}

// loads all reminders from the db and schedules them.
func StartAllReminders(s *discordgo.Session) error {
	reminders, err := repository.ReminderService.DBClient.GetAllReminders()
	if err != nil {
		return fmt.Errorf("failed to load reminders: %v", err)
	}
	for _, r := range reminders {
		err = StartReminder(s, r.UserID, r.Message, r.Email, r.Hour, r.Minute, r.AM)
		if err != nil {
			fmt.Printf("failed to start reminder for user %v: %v\n", r.UserID, err)
		}
	}
	return nil
}