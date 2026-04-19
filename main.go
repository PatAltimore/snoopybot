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

	// Azure storage is always required.
	for _, k := range []string{"AZURE_STORAGE_ACCOUNT", "AZURE_STORAGE_ACCESS_KEY"} {
		if os.Getenv(k) == "" {
			log.Fatalf("missing required env var: %s", k)
		}
	}

	// At least one posting platform must be configured.
	mastodonOK := os.Getenv("MASTODON_SERVER") != "" && os.Getenv("MASTODON_ACCESS_TOKEN") != ""
	// threadsOK := os.Getenv("THREADS_USER_ID") != "" && os.Getenv("THREADS_ACCESS_TOKEN") != "" // re-enable for Threads
	if !mastodonOK { // add `&& !threadsOK` when Threads is re-enabled
		log.Fatal("at least one posting platform must be configured: " +
			"set MASTODON_SERVER + MASTODON_ACCESS_TOKEN")
	}

	log.Println("Snoopy is writing...")
	if err := bot.DoWork(); err != nil {
		log.Fatalf("bot error: %v", err)
	}
	log.Println("Snoopy is done.")
}
