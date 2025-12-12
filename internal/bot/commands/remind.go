package commands

import (
	"fmt"

	"github.com/Sush1sui/meds_reminder/internal/common"
	"github.com/bwmarrin/discordgo"
)

func RemindCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	// Get the user choice
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please select a user to remind.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	userChoice := options[0].StringValue()

	// Validate user choice
	if userChoice != "jp" && userChoice != "dane" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid selection. Please choose JP or Dane.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
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
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Remind command panic: %v\n", r)
			}
		}()

		// send reminder based on choice
		if userChoice == "jp" {
			common.RemindUser(s)
		} else if userChoice == "dane" {
			common.SendSimpleReminder(s)
		}
	}()

	// final response
	var responseMsg string
	if userChoice == "jp" {
		responseMsg = "JP medication reminder sent. Check logs for details."
	} else {
		responseMsg = "Dane simple reminder sent. Check logs for details."
	}

	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: func(s string) *string { return &s }(responseMsg),
	})
}