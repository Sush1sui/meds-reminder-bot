package commands

import (
	"strconv"

	"github.com/Sush1sui/meds_reminder/internal/common"
	"github.com/Sush1sui/meds_reminder/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func GetUserReminders(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Fetching reminders...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	user := i.ApplicationCommandData().GetOption("user").UserValue(s)
	if user == nil {
		common.ResponseEdit(s, i, "Missing required user option.")
		return
	}

	reminders, err := repository.ReminderService.DBClient.GetUserReminders(user.ID)
	if err != nil {
		common.ResponseEdit(s, i, "Failed to fetch reminders.")
		return
	}
	if len(reminders) == 0 {
		common.ResponseEdit(s, i, "No reminders found for the user.")
		return
	}

	response := "Reminders for " + user.Username + ":\n"
	for i := range reminders {
		r := reminders[i]

		response += "```"
		response += "ID: " + r.ID.Hex() + "\n"
		response += "Message: " + r.Message + "\n"
		response += "Time: " + strconv.Itoa(r.Hour) + ":"
		if r.Minute == 0 {
			response += "00 "
		} else {
			response += strconv.Itoa(r.Minute) + " "
		}
		if r.AM {
			response += "AM" + "\n"
		} else {
			response += "PM" + "\n"
		}
		response += "```"
	}

	common.ResponseEdit(s, i, response)
}