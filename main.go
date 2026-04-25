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

	for _, k := range []string{"AZURE_STORAGE_ACCOUNT", "AZURE_STORAGE_ACCESS_KEY"} {
		if os.Getenv(k) == "" {
			log.Fatalf("missing required env var: %s", k)
		}
	}

	if os.Getenv("MASTODON_SERVER") == "" || os.Getenv("MASTODON_ACCESS_TOKEN") == "" {
		log.Fatal("MASTODON_SERVER and MASTODON_ACCESS_TOKEN must both be set")
	}

	log.Println("Snoopy is writing...")
	if err := bot.DoWork(); err != nil {
		log.Fatalf("bot error: %v", err)
	}
	log.Println("Snoopy is done.")
}
