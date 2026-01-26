package mongodb

import (
	"context"
	"fmt"

	"github.com/Sush1sui/meds_reminder/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (c *MongoClient) GetAllReminders() ([]*model.Reminder, error) {
	var reminders []*model.Reminder

	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil { return nil, fmt.Errorf("Error fetching reminders: %v", err) }
	
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var reminder model.Reminder
		if err := cursor.Decode(&reminder); err != nil {
			return nil, fmt.Errorf("Error decoding reminder: %v", err)
		}
		reminders = append(reminders, &reminder)
	}
	return reminders, nil
}

func (c *MongoClient) GetReminder(userID, message, email string, hour, minute int, am bool) (*model.Reminder, error) {
	filter := bson.M{
		"user_id": userID,
		"message": message,
		"email":   email,
		"hour":    hour,
		"minute":  minute,
		"am":      am,
	}
	var reminder model.Reminder
	err := c.Client.FindOne(context.Background(), filter).Decode(&reminder)
	if err != nil {
		return nil, fmt.Errorf("Error fetching reminder: %v", err)
	}
	return &reminder, nil
}

func (c *MongoClient) CreateReminder(userID, message, email string, hour, minute int, am bool) (*model.Reminder, error) {
	if userID == "" || message == "" {
		return nil, fmt.Errorf("UserID and Message cannot be empty")
	}
	if hour < 1 || hour > 12 || minute < 0 || minute > 59 {
		return nil, fmt.Errorf("Invalid time format")
	}
	reminder := &model.Reminder{
		UserID:  userID,
		Message: message,
		Email:   email,
		Hour:    hour,
		Minute:  minute,
		AM:      am,
	}

	result, err := c.Client.InsertOne(context.Background(), reminder)
	if err != nil { return nil, fmt.Errorf("Error inserting reminder: %v", err) }

	reminder.ID = result.InsertedID.(bson.ObjectID)
	return reminder, nil
}

func (c *MongoClient) GetUserReminders(userID string) ([]*model.Reminder, error) {
	if userID == "" { return nil, fmt.Errorf("UserID cannot be empty") }

	var reminders []*model.Reminder
	cursor, err := c.Client.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil { return nil, fmt.Errorf("Error fetching user reminders: %v", err) }

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var reminder model.Reminder
		if err := cursor.Decode(&reminder); err != nil {
			return nil, fmt.Errorf("Error decoding reminder: %v", err)
		}
		reminders = append(reminders, &reminder)
	}

	return reminders, nil
}

func (c *MongoClient) DeleteReminder(reminderID string) error {
	objID, err := bson.ObjectIDFromHex(reminderID)
	if err != nil { return fmt.Errorf("Invalid reminder ID: %v", err) }
	result, err := c.Client.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil { return fmt.Errorf("Error deleting reminder: %v", err) }
	if result.DeletedCount == 0 {
		return fmt.Errorf("No reminder found with the given ID")
	}
	return nil
}