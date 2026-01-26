package commands

import (
	"github.com/Sush1sui/meds_reminder/internal/common"
	"github.com/Sush1sui/meds_reminder/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func DeleteReminder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Deleteing Reminder...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	reminderID := i.ApplicationCommandData().GetOption("reminder_id").StringValue()
	if reminderID == "" {
		common.ResponseEdit(s, i, "Missing required reminder_id option.")
		return
	}

	err := repository.ReminderService.DBClient.DeleteReminder(reminderID)
	if err != nil {
		common.ResponseEdit(s, i, "Something went wrong with deleting the reminder.")
		return
	}

	common.ResponseEdit(s, i, "Reminder deleted successfully.")
}