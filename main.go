package main

import (
	"log"
	"os"

	"github.com/PatAltimore/snoopybot/internal/bot"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present. Silently ignored in production where no .env exists.
	_ = godotenv.Load()

	required := []string{
		"TWITTER_CONSUMER_KEY",
		"TWITTER_CONSUMER_SECRET",
		"TWITTER_ACCESS_TOKEN",
		"TWITTER_ACCESS_TOKEN_SECRET",
		"AZURE_STORAGE_ACCOUNT",
		"AZURE_STORAGE_ACCESS_KEY",
	}
	for _, k := range required {
		if os.Getenv(k) == "" {
			log.Fatalf("missing required env var: %s", k)
		}
	}

	log.Println("Snoopy is writing...")
	if err := bot.DoWork(); err != nil {
		log.Fatalf("bot error: %v", err)
	}
	log.Println("Snoopy is done.")
}
