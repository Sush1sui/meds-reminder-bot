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
		log.Fatal("Bot token not found in environment variables")
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

	// Start simple reminder for Dane at 10am
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Simple scheduler panic: %v\n", r)
			}
		}()

		loc := time.FixedZone("UTC+8", 8*3600)
		for {
			now := time.Now().In(loc)
			
			// Check if it's 10:00 AM
			if now.Hour() == 10 && now.Minute() == 0 {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("Simple reminder panic: %v\n", r)
						}
					}()
					common.SendSimpleReminder(sess)
				}()
				
					// Sleep for 61 seconds to avoid duplicate in same minute
				time.Sleep(61 * time.Second)
			} else {
				// Sleep until next minute boundary
				nextMinute := now.Truncate(time.Minute).Add(time.Minute)
				sleepDuration := time.Until(nextMinute)
				if sleepDuration > 0 {
					time.Sleep(sleepDuration)
				}
			}
		}
	}()

	// Medication reminder scheduler for JP has been removed (medications inactive)

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

