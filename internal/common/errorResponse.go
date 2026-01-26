package common

import "github.com/bwmarrin/discordgo"

func ResponseEdit(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: func(s string) *string { return &s }(message),
	})
}