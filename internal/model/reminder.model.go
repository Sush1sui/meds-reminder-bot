package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Reminder struct {
	ID 			bson.ObjectID `bson:"_id,omitempty"`
	UserID	 	string        `bson:"user_id"`
	Message   	string        `bson:"message"`
	Email 	 	string        `bson:"email"`
	Hour    	int           `bson:"hour"`
	Minute 	 	int           `bson:"minute"`
	AM 			bool          `bson:"am"`
}