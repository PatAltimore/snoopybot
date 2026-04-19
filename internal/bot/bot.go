package bot

import (
	"context"
	"log"
	"math/rand"
	"os"

	"github.com/PatAltimore/snoopybot/internal/mastodon"
	"github.com/PatAltimore/snoopybot/internal/storage"
)

func DoWork() error {
	ctx := context.Background()

	mastodonClient := &mastodon.Client{
		Server:      os.Getenv("MASTODON_SERVER"),
		AccessToken: os.Getenv("MASTODON_ACCESS_TOKEN"),
	}

	stateClient, err := storage.NewStateClient(
		os.Getenv("AZURE_STORAGE_ACCOUNT"),
		os.Getenv("AZURE_STORAGE_ACCESS_KEY"),
	)
	if err != nil {
		return err
	}

	if err := stateClient.EnsureTable(ctx); err != nil {
		return err
	}

	dryRun := os.Getenv("DRY_RUN") == "true"

	switch rand.Intn(2) {
	case 0:
		index := stateClient.GetNovelIndex(ctx)
		index = (index + 1) % len(novel)
		if err := stateClient.SetNovelIndex(ctx, index); err != nil {
			return err
		}
		log.Printf("Posting novel index: %d", index)
		return postStatus(ctx, mastodonClient, novel[index], dryRun)

	case 1:
		phrase := rand.Intn(len(misc))
		log.Printf("Posting miscellaneous quote index: %d", phrase)
		return postStatus(ctx, mastodonClient, misc[phrase], dryRun)
	}
	return nil
}

func postStatus(ctx context.Context, c *mastodon.Client, text string, dryRun bool) error {
	log.Printf("Status: %s", text)
	if dryRun {
		log.Println("DRY_RUN=true, skipping actual post")
		return nil
	}
	return c.PostStatus(ctx, text)
}
