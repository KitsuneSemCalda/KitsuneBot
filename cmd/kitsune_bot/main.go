package main

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"KitsuneSemCalda/KitsuneBot/internal/twitch"
	"log"
)

func main() {
	config := twitch.LoadConfig()
	client := twitch.NewClient(config)
	registry := twitch.NewHandlerRegistry()
	err := db.Setup("KitsuneBot.db")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = db.CleanupTTL(db.DB, 30)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	defer db.DB.Close()

	registry.Register(twitch.NewPingHandler())
	registry.Register(twitch.NewStatsHandler())
	registry.Register(twitch.NewHelpHandler(registry))

	client.RegisterHandlers(registry)

	log.Printf("Starting KitsuneBot for channel: %s", config.Channel)
	if err := client.Connect(); err != nil {
		log.Fatalf("Error connecting to Twitch: %v", err)
	}

}
