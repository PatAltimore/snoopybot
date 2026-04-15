package bot

import (
	"context"
	"log"
	"math/rand"
	"os"

	"github.com/PatAltimore/snoopybot/internal/storage"
	"github.com/PatAltimore/snoopybot/internal/twitter"
)

func DoWork() error {
	ctx := context.Background()

	twitterClient := &twitter.Client{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		AccessToken:    os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
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
		log.Printf("Tweeting novel index: %d", index)
		return sendTweet(ctx, twitterClient, novel[index], dryRun)

	case 1:
		phrase := rand.Intn(len(misc))
		log.Printf("Tweeting miscellaneous quote index: %d", phrase)
		return sendTweet(ctx, twitterClient, misc[phrase], dryRun)
	}
	return nil
}

func sendTweet(ctx context.Context, c *twitter.Client, text string, dryRun bool) error {
	log.Printf("Tweet: %s", text)
	if dryRun {
		log.Println("DRY_RUN=true, skipping actual tweet")
		return nil
	}
	return c.PostTweet(ctx, text)
}
