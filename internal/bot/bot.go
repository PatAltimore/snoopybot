package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/PatAltimore/snoopybot/internal/mastodon"
	"github.com/PatAltimore/snoopybot/internal/storage"
	// "github.com/PatAltimore/snoopybot/internal/threads" // re-enable to post to Threads
)

// poster is satisfied by any platform client (Mastodon, Threads, etc.)
type poster interface {
	PostStatus(ctx context.Context, text string) error
	Name() string
}

func DoWork() error {
	ctx := context.Background()
	dryRun := os.Getenv("DRY_RUN") == "true"

	// Build the list of configured platforms — any combination is valid.
	var posters []poster
	if server := os.Getenv("MASTODON_SERVER"); server != "" {
		posters = append(posters, &mastodon.Client{
			Server:      server,
			AccessToken: os.Getenv("MASTODON_ACCESS_TOKEN"),
		})
	}
	// Threads — uncomment to re-enable:
	// if userID := os.Getenv("THREADS_USER_ID"); userID != "" {
	// 	posters = append(posters, &threads.Client{
	// 		UserID:      userID,
	// 		AccessToken: os.Getenv("THREADS_ACCESS_TOKEN"),
	// 	})
	// }

	// Select content
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

	var text string
	switch rand.Intn(2) {
	case 0:
		index := stateClient.GetNovelIndex(ctx)
		index = (index + 1) % len(novel)
		if err := stateClient.SetNovelIndex(ctx, index); err != nil {
			return err
		}
		log.Printf("Selected novel index: %d", index)
		text = novel[index]
	case 1:
		phrase := rand.Intn(len(misc))
		log.Printf("Selected miscellaneous quote index: %d", phrase)
		text = misc[phrase]
	}

	log.Printf("Status: %s", text)

	// Post to every configured platform; collect all errors.
	var errs []error
	for _, p := range posters {
		if dryRun {
			log.Printf("DRY_RUN=true, skipping post to %s", p.Name())
			continue
		}
		if err := p.PostStatus(ctx, text); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", p.Name(), err))
		} else {
			log.Printf("Posted to %s", p.Name())
		}
	}
	return errors.Join(errs...)
}
