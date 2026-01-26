package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Config holds the application's configuration
type Config struct {
	DiscordToken string
	ServerPort   string
	AppID        string
	ServerURL    string
}

var GlobalConfig Config

// LoadConfig initializes the configuration with default values
func LoadConfig() (error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	GlobalConfig = Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		ServerPort:   os.Getenv("SERVER_PORT"),
		AppID:        os.Getenv("APP_ID"),
		ServerURL:    os.Getenv("SERVER_URL"),
	}
	return nil
}

func MongoConnection() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("DB Connected!")

	return client
}