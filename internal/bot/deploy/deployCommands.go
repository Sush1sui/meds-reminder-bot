package deploy

import (
	"fmt"
	"log"

	"github.com/Sush1sui/meds_reminder/internal/bot/commands"
	"github.com/bwmarrin/discordgo"
)

// List all slash commands here
var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "add-reminder",
		Description: "Add a reminder",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Choose who to remind",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:		"email",
				Description:	"User email",
				Required:	true,
			},
			{
				Type:		discordgo.ApplicationCommandOptionString,
				Name:		"message",
				Description:	"Reminder message",
				Required:	true,
			},
			{
				Type:		discordgo.ApplicationCommandOptionInteger,
				Name:		"hour",
				Description:	"Hour (1-12)",
				Required:	true,
			},
			{
				Type:		discordgo.ApplicationCommandOptionInteger,
				Name:		"minute",
				Description:	"Minute (0-59)",
				Required:	true,
			},
			{
				Type:		discordgo.ApplicationCommandOptionBoolean,
				Name:		"am",
				Description:	"True for AM, False for PM",
				Required:	true,
			},
		},
	},
	{
		Name:		"get-all-reminders",
		Description:	"Get your reminders",
	},
	{
		Name: 		"get-user-reminders",
		Description: 	"Get reminders for a specific user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Select the user",
				Required:    true,
			},
		},
	},
	{
		Name:		"delete-reminder",
		Description:	"Delete a reminder by its ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:		discordgo.ApplicationCommandOptionString,
				Name:		"reminder_id",
				Description:	"ID of the reminder to delete",
				Required:	true,
			},
		},
	},
	// Add more commands here
}

// Map command names to handler functions
var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"add-reminder":		commands.AddReminder,
	"get-all-reminders":		commands.GetAllReminders,
	"get-user-reminders":	commands.GetUserReminders,
	"delete-reminder":	commands.DeleteReminder,
	// Add more: "hello": commands.HelloCommand, etc.
}

func DeployCommands(sess *discordgo.Session) {
	// Remove all global commands
	globalCmds, err := sess.ApplicationCommands(sess.State.User.ID, "")
	if err == nil {
			for _, cmd := range globalCmds {
					err := sess.ApplicationCommandDelete(sess.State.User.ID, "", cmd.ID)
					if err != nil {
							log.Printf("Failed to delete global command %s: %v", cmd.Name, err)
					} else {
							log.Printf("Deleted global command: %s", cmd.Name)
					}
			}
	}

	// Bulk overwrite commands for each guild (this replaces all commands)
	guilds := sess.State.Guilds
	for _, guild := range guilds {
			_, err := sess.ApplicationCommandBulkOverwrite(sess.State.User.ID, guild.ID, SlashCommands)
			if err != nil {
					log.Fatalf("Cannot create slash commands for guild %s: %v", guild.ID, err)
			}
	}

	// Register handler for slash commands
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if handler, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
					handler(s, i)
			} else {
					fmt.Printf("Unknown command: %s\n", i.ApplicationCommandData().Name)
					fmt.Printf("Available commands: %v\n", CommandHandlers)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
									Content: "Unknown command.",
									Flags:   discordgo.MessageFlagsEphemeral,
							},
					})
			}
	})

	log.Println("Slash commands deployed successfully.")
}