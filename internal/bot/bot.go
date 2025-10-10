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

    go func() {
        loc := time.FixedZone("UTC+8", 8*3600)
        for {
            now := time.Now().In(loc)
            // next 10:00 in UTC+8
            next := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, loc)
            if !next.After(now) {
                next = next.Add(24 * time.Hour)
            }
            sleepDur := time.Until(next)
            // sleep until the next scheduled time
            time.Sleep(sleepDur)

            // call reminder (recover from panics to keep loop running)
            func() {
                defer func() {
                    if r := recover(); r != nil {
                        fmt.Printf("reminder panic: %v\n", r)
                    }
                }()
                common.RemindUser(sess)
            }()
        }
    }()

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

