package common

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Medication represents a single medication with its schedule
type Medication struct {
	Name           string   `json:"name"`
	Dose           string   `json:"dose"`
	Times          []string `json:"times"` // e.g., ["06:00", "18:00"]
	DaysRemaining  int      `json:"days_remaining"`
	TotalDays      int      `json:"total_days"`
	Indication     string   `json:"indication"`
	Notes          string   `json:"notes,omitempty"`
	Active         bool     `json:"active"`
	StartDate      string   `json:"start_date"` // YYYY-MM-DD format
}

// MedicationSchedule holds all medications and their states
type MedicationSchedule struct {
	Medications      []Medication `json:"medications"`
	LastUpdated      string       `json:"last_updated"`
	PrednisoneSwitch bool         `json:"prednisone_switch"` // true when switched to 1x/day
}

const stateFile = "medication_state.json"

var stateMutex sync.Mutex

// LoadMedicationState loads the medication state from file
func LoadMedicationState() (*MedicationSchedule, error) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize with default schedule if file doesn't exist
			// Unlock before calling init to avoid deadlock
			stateMutex.Unlock()
			schedule := initializeDefaultSchedule()
			stateMutex.Lock()
			return schedule, nil
		}
		return nil, err
	}

	var schedule MedicationSchedule
	if err := json.Unmarshal(data, &schedule); err != nil {
		return nil, err
	}

	return &schedule, nil
}

// SaveMedicationState saves the medication state to file
func SaveMedicationState(schedule *MedicationSchedule) error {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	schedule.LastUpdated = time.Now().Format(time.RFC3339)
	data, err := json.MarshalIndent(schedule, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(stateFile, data, 0644)
}

// Hardcoded user credentials for JP's medication reminders
const (
	JPDiscordID = "982491279369830460"
	JPEmails    = "jpmercado57123@gmail.com,kingbluetoothz24@gmail.com,jplovesraw@gmail.com,jpmercado57@yahoo.com"
)

// initializeDefaultSchedule creates the initial medication schedule
func initializeDefaultSchedule() *MedicationSchedule {
	startDate := "2025-12-11" // Started December 11, 2025 at 6pm

	schedule := &MedicationSchedule{
		LastUpdated:      time.Now().Format(time.RFC3339),
		PrednisoneSwitch: false,
		Medications: []Medication{
			{
				Name:          "Doxycycline Hydrochloride 50mg/ml (EHRLICURE)",
				Dose:          "2 ml",
				Times:         []string{"18:00"}, // 6pm
				DaysRemaining: 28,
				TotalDays:     28,
				Indication:    "Antibacterial",
				Active:        true,
				StartDate:     startDate,
			},
			{
				Name:          "Prednisone 20mg tab",
				Dose:          "1/2 tab",
				Times:         []string{"09:00", "21:00"}, // 9am, 9pm
				DaysRemaining: 7,
				TotalDays:     7,
				Indication:    "Corticosteroid",
				Notes:         "After 7 days, switch to 9pm only",
				Active:        true,
				StartDate:     startDate,
			},
			{
				Name:          "Papi Bion Plus",
				Dose:          "1.5 ml",
				Times:         []string{"13:00"}, // 1pm
				DaysRemaining: 30,
				TotalDays:     30,
				Indication:    "Blood supplement",
				Notes:         "Don't combine with Doxycycline on short intervals",
				Active:        true,
				StartDate:     startDate,
			},
			{
				Name:          "Thromb Beat",
				Dose:          "4 ml",
				Times:         []string{"08:00", "20:00"}, // 8am, 8pm
				DaysRemaining: 30,
				TotalDays:     30,
				Indication:    "Platelet supplement",
				Active:        true,
				StartDate:     startDate,
			},
			{
				Name:          "Immunol syrup",
				Dose:          "2 ml",
				Times:         []string{"07:00", "19:00"}, // 7am, 7pm
				DaysRemaining: 30,
				TotalDays:     30,
				Indication:    "Immune booster",
				Active:        true,
				StartDate:     startDate,
			},
			{
				Name:          "Livertel",
				Dose:          "3 ml",
				Times:         []string{"07:15", "19:15"}, // 7:15am, 7:15pm
				DaysRemaining: 15,
				TotalDays:     15,
				Indication:    "Liver supplement",
				Active:        true,
				StartDate:     startDate,
			},
			{
				Name:          "Sync Nephric syrup",
				Dose:          "3 ml",
				Times:         []string{"07:30", "19:30"}, // 7:30am, 7:30pm
				DaysRemaining: 15,
				TotalDays:     15,
				Indication:    "Kidney supplement",
				Active:        true,
				StartDate:     startDate,
			},
		},
	}

	// Save the initial state
	if err := SaveMedicationState(schedule); err != nil {
		fmt.Printf("Warning: failed to save initial medication state: %v\n", err)
	}

	return schedule
}

// UpdateMedicationCounts updates the day counts based on elapsed days
func UpdateMedicationCounts(schedule *MedicationSchedule) {
	loc := time.FixedZone("UTC+8", 8*3600)
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")

	// Check if we already updated today
	if schedule.LastUpdated != "" {
		lastUpdate, err := time.Parse(time.RFC3339, schedule.LastUpdated)
		if err == nil {
			lastUpdateDay := lastUpdate.In(loc).Format("2006-01-02")
			if lastUpdateDay == today {
				return // Already updated today
			}
		}
	}

	for i := range schedule.Medications {
		med := &schedule.Medications[i]
		if !med.Active {
			continue
		}

		startDate, err := time.Parse("2006-01-02", med.StartDate)
		if err != nil {
			fmt.Printf("Error parsing start date for %s: %v\n", med.Name, err)
			continue
		}

		// Calculate days elapsed since start
		daysElapsed := int(now.Sub(startDate).Hours() / 24)
		
		// Update days remaining
		med.DaysRemaining = med.TotalDays - daysElapsed
		
		if med.DaysRemaining <= 0 {
			med.Active = false
			med.DaysRemaining = 0
		}

		// Special handling for Prednisone: switch to 1x/day after 7 days
		if med.Name == "Prednisone 20mg tab" && daysElapsed >= 7 && !schedule.PrednisoneSwitch {
			med.Times = []string{"21:00"} // Switch to 9pm only
			schedule.PrednisoneSwitch = true
			fmt.Println("Prednisone switched to 1x per day (9pm only)")
		}
	}

	SaveMedicationState(schedule)
}

// GetCurrentReminders returns the list of medications due at or around the given time
func GetCurrentReminders(schedule *MedicationSchedule, currentTime time.Time) []Medication {
	loc := time.FixedZone("UTC+8", 8*3600)
	now := currentTime.In(loc)
	currentHour := now.Format("15:04")

	var reminders []Medication
	for _, med := range schedule.Medications {
		if !med.Active {
			continue
		}

		for _, medTime := range med.Times {
			if medTime == currentHour {
				reminders = append(reminders, med)
				break
			}
		}
	}

	return reminders
}

// FormatReminderMessage creates a formatted reminder message
func FormatReminderMessage(reminders []Medication) string {
	if len(reminders) == 0 {
		return ""
	}

	message := "🔔 **Medication Reminder** 🔔\n\n"
	message += "It's time for Diluc's meds! 💊✨\n\n"

	for _, med := range reminders {
		message += fmt.Sprintf("💊 **%s**\n", med.Name)
		message += fmt.Sprintf("   📏 Dose: %s\n", med.Dose)
		message += fmt.Sprintf("   🏥 Purpose: %s\n", med.Indication)
		message += fmt.Sprintf("   📅 Days remaining: %d/%d\n", med.DaysRemaining, med.TotalDays)
		if med.Notes != "" {
			message += fmt.Sprintf("   ℹ️ Note: %s\n", med.Notes)
		}
		message += "\n"
	}

	message += "Remember to take with/after meals! 🍽️\n"
	message += "Stay strong! 💪✨"

	return message
}

// GetAllActiveMedications returns a summary of all active medications
func GetAllActiveMedications(schedule *MedicationSchedule) string {
	message := "📋 **Current Medication Schedule** 📋\n\n"

	for _, med := range schedule.Medications {
		if !med.Active {
			continue
		}

		timesStr := ""
		for i, t := range med.Times {
			if i > 0 {
				timesStr += ", "
			}
			timesStr += t
		}

		message += fmt.Sprintf("💊 **%s**\n", med.Name)
		message += fmt.Sprintf("   📏 Dose: %s\n", med.Dose)
		message += fmt.Sprintf("   ⏰ Times: %s\n", timesStr)
		message += fmt.Sprintf("   🏥 Purpose: %s\n", med.Indication)
		message += fmt.Sprintf("   📅 Days remaining: %d/%d\n", med.DaysRemaining, med.TotalDays)
		if med.Notes != "" {
			message += fmt.Sprintf("   ℹ️ Note: %s\n", med.Notes)
		}
		message += "\n"
	}

	return message
}
