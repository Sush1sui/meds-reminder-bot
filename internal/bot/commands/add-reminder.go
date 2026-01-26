package commands

import (
	"strconv"

	"github.com/Sush1sui/meds_reminder/internal/common"
	"github.com/Sush1sui/meds_reminder/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func AddReminder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
			Content: "Adding reminder...",
            Flags: discordgo.MessageFlagsEphemeral,
        },
    })
	if err != nil {
		common.ResponseEdit(s, i, "Failed to respond to interaction.")
		return
	}

	// Get user input from the interaction
	user := i.ApplicationCommandData().GetOption("user").UserValue(s)
	email := i.ApplicationCommandData().GetOption("email").StringValue()
	message := i.ApplicationCommandData().GetOption("message").StringValue()
	hour := i.ApplicationCommandData().GetOption("hour").IntValue()
	minute := i.ApplicationCommandData().GetOption("minute").IntValue()
	am := i.ApplicationCommandData().GetOption("am").BoolValue()

	if user == nil || email == "" || message == "" || hour < 1 || hour > 12 || minute < 0 || minute > 59 {
		common.ResponseEdit(s, i, "Missing required options.")
		return
	}

	existingReminder, err := repository.ReminderService.DBClient.GetReminder(
		user.ID,
		message,
		email,
		int(hour),
		int(minute),
		am,
	)
	if err == nil && existingReminder != nil {
		common.ResponseEdit(s, i, "Reminder already exists.")
		return
	}

	reminder, err := repository.ReminderService.DBClient.CreateReminder(
		user.ID,
		message,
		email,
		int(hour),
		int(minute),
		am,
	)

	if err != nil {
		common.ResponseEdit(s, i, "Failed to create reminder.")
		return
	}

	amOrPm := "AM"
	if !am { amOrPm = "PM" }
	messageResponse := "Reminder added successfully! ID: " + reminder.ID.String() + "\n" +
	"User: " + user.Username + "\n" +
	"Email: " + email + "\n" +
	"Message: " + message + "\n" +
	"Hour: " + strconv.FormatInt(hour, 10) + "\n" +
	"Minute: " + strconv.FormatInt(minute, 10) + "\n" +
	"AM/PM: " + amOrPm

	common.ResponseEdit(s, i, messageResponse)
}