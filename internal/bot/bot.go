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
		// run every 10am at UTC+8
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			now := time.Now().In(time.FixedZone("UTC+8", 8*3600))
			if now.Hour() == 10 && now.Minute() == 0 {
				common.RemindUser(sess)
			}
			<-ticker.C
		}
	}()

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

