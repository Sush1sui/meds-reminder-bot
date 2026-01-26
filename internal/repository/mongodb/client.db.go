package mongodb

import "go.mongodb.org/mongo-driver/v2/mongo"

type MongoClient struct {
	Client *mongo.Collection
}