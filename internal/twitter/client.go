package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dghubble/oauth1"
)

const tweetsURL = "https://api.twitter.com/2/tweets"

type Client struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func (c *Client) PostTweet(ctx context.Context, text string) error {
	cfg := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
	httpClient := cfg.Client(ctx, token)

	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return err
	}

	resp, err := httpClient.Post(tweetsURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("twitter API returned status %d", resp.StatusCode)
	}
	return nil
}
