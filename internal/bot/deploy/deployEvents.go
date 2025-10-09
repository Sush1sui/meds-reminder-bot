package deploy

import (
	"log"

	"github.com/bwmarrin/discordgo"
)


var EventHandlers = []any{
	// Add more event handlers here, e.g.:
	// Go doesn't support dynamic runtime imports
	// You have to manually add each event handler
}

func DeployEvents(sess *discordgo.Session) {
	for _, handler := range EventHandlers {
		sess.AddHandler(handler)
	}
	log.Println("Event handlers deployed successfully.")
}