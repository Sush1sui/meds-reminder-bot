package repository

import "github.com/Sush1sui/meds_reminder/internal/model"

type ReminderInterface interface {
	GetAllReminders() ([]*model.Reminder, error)
	CreateReminder(string, string, string, int, int, bool) (*model.Reminder, error)
	GetReminder(string, string, string, int, int, bool) (*model.Reminder, error)
	GetUserReminders(string) ([]*model.Reminder, error)
	DeleteReminder(string) error
}

type ReminderServiceType struct {
	DBClient ReminderInterface
}

var ReminderService ReminderServiceType