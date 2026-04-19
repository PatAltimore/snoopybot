package threads

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "https://graph.threads.net/v1.0"

// Client posts to Threads using Meta's two-step Graph API:
// 1. Create a text container  →  POST /{user-id}/threads
// 2. Publish the container    →  POST /{user-id}/threads_publish
//
// Obtain a long-lived access token with threads_basic and
// threads_content_publish permissions from the Meta developer portal.
type Client struct {
	UserID      string
	AccessToken string
}

func (c *Client) Name() string { return "Threads" }

func (c *Client) PostStatus(ctx context.Context, text string) error {
	containerID, err := c.createContainer(ctx, text)
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}
	return c.publish(ctx, containerID)
}

func (c *Client) createContainer(ctx context.Context, text string) (string, error) {
	params := url.Values{}
	params.Set("media_type", "TEXT")
	params.Set("text", text)
	params.Set("access_token", c.AccessToken)

	endpoint := fmt.Sprintf("%s/%s/threads?%s", baseURL, c.UserID, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("threads API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing container response: %w", err)
	}
	if result.ID == "" {
		return "", fmt.Errorf("threads API returned empty container ID: %s", string(body))
	}
	return result.ID, nil
}

func (c *Client) publish(ctx context.Context, containerID string) error {
	params := url.Values{}
	params.Set("creation_id", containerID)
	params.Set("access_token", c.AccessToken)

	endpoint := fmt.Sprintf("%s/%s/threads_publish?%s", baseURL, c.UserID, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("threads publish returned status %d: %s", resp.StatusCode, string(detail))
	}
	return nil
}
