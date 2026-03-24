package main

import (
	"KitsuneSemCalda/KitsuneBot/internal/db"
	"log"
)

func main() {
	err := db.Setup("KitsuneBot.db")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = db.CleanupTTL(db.DB, 30)

	if err != nil {
		log.Printf("Error: %v", err)
	}

}
