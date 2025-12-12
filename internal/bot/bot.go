package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sush1sui/meds_reminder/internal/bot/deploy"
	"github.com/Sush1sui/meds_reminder/internal/common"
	"github.com/Sush1sui/meds_reminder/internal/config"
	"github.com/bwmarrin/discordgo"
)

func StartBot() {

	// create new discord session
	if config.GlobalConfig.DiscordToken == "" {
		fmt.Println("Bot token not found")
	}
	sess, err := discordgo.New("Bot " + config.GlobalConfig.DiscordToken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildPresences | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
    s.UpdateStatusComplex(discordgo.UpdateStatusData{
        Status: "idle",
    })
	})

	err = sess.Open()
	if err != nil {
		log.Fatalf("error opening connection to Discord: %v", err)
	}
	defer sess.Close()

	// Deploy commands
	deploy.DeployCommands(sess)

	// Deploy events
	deploy.DeployEvents(sess)

	// Start medication reminder scheduler
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Scheduler panic: %v\n", r)
			}
		}()

		loc := time.FixedZone("UTC+8", 8*3600)
		
		// Pre-parse reminder times for performance
		type reminderTime struct {
			hour   int
			minute int
		}
		
		reminderTimesList := []reminderTime{
			{7, 0}, {7, 15}, {7, 30},  // Morning meds
			{8, 0},                     // Thromb Beat morning
			{9, 0},                     // Prednisone morning (first 7 days)
			{13, 0},                    // Papi Bion Plus
			{18, 0},                    // Doxycycline
			{19, 0}, {19, 15}, {19, 30}, // Evening meds
			{20, 0},                    // Thromb Beat evening
			{21, 0},                    // Prednisone evening
		}

		lastTriggered := time.Time{}

		for {
			now := time.Now().In(loc)
			currentTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, loc)

			// Avoid duplicate triggers in the same minute
			if currentTime.Equal(lastTriggered) {
				time.Sleep(30 * time.Second)
				continue
			}

			// Check if current time matches any reminder time
			for _, rt := range reminderTimesList {
				if now.Hour() == rt.hour && now.Minute() == rt.minute {
					lastTriggered = currentTime
					
					// Call reminder in goroutine to avoid blocking scheduler
					go func() {
						defer func() {
							if r := recover(); r != nil {
								fmt.Printf("reminder panic: %v\n", r)
							}
						}()
						common.RemindUser(sess)
					}()
					
					break
				}
			}

			// Sleep until the next minute boundary
			nextMinute := now.Truncate(time.Minute).Add(time.Minute)
			sleepDuration := time.Until(nextMinute)
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}
		}
	}()

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

