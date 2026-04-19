package mastodon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client posts statuses to a Mastodon instance using a Bearer access token.
// The token is created in your instance's developer settings with the
// write:statuses scope.
type Client struct {
	Server      string // e.g. "https://mastodon.social"
	AccessToken string
}

func (c *Client) PostStatus(ctx context.Context, text string) error {
	url := c.Server + "/api/v1/statuses"

	body, err := json.Marshal(map[string]string{"status": text})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("mastodon API returned status %d: %s", resp.StatusCode, string(detail))
	}
	return nil
}
