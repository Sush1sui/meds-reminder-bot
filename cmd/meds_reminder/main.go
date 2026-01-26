package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/meds_reminder/internal/bot"
	"github.com/Sush1sui/meds_reminder/internal/config"
	"github.com/Sush1sui/meds_reminder/internal/repository"
	"github.com/Sush1sui/meds_reminder/internal/repository/mongodb"
	"github.com/Sush1sui/meds_reminder/internal/server/routes"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	mongoClient := config.MongoConnection()
	defer mongoClient.Disconnect(context.Background())
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	remindersCollection := mongoClient.Database(os.Getenv("MONGODB_NAME")).Collection("reminders")

	repository.ReminderService = repository.ReminderServiceType{
		DBClient: &mongodb.MongoClient{
			Client: remindersCollection,
		},
	}

	addr := fmt.Sprintf(":%s", config.GlobalConfig.ServerPort)
	router := routes.NewRouter()
	fmt.Printf("Server listening on Port:%s\n", config.GlobalConfig.ServerPort)

	// Run HTTP server in a goroutine
	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			// Log error instead of panicking to avoid crashing the service
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Run Discord bot in a goroutine
	go bot.StartBot()
	// Block main goroutine until interrupt signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down...")
}