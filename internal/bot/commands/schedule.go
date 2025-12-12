package commands

import (
	"github.com/Sush1sui/meds_reminder/internal/common"
	"github.com/bwmarrin/discordgo"
)

func ScheduleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Load medication schedule
	schedule, err := common.LoadMedicationState()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error loading medication schedule: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Update counts
	common.UpdateMedicationCounts(schedule)

	// Get formatted schedule
	scheduleMsg := common.GetAllActiveMedications(schedule)

	// Respond with the schedule
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: scheduleMsg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
