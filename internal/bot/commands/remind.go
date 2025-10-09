package commands

import (
	"github.com/Sush1sui/meds_reminder/internal/common"
	"github.com/bwmarrin/discordgo"
)

func RemindCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	// acknowledge immediately so we have time to do work
    _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Flags: discordgo.MessageFlagsEphemeral,
        },
    })

	go func() {
		// send reminder
		common.RemindUser(s)
	}()

	// final response
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: func(s string) *string { return &s }("Reminder process started. Check logs for details."),
	})
}