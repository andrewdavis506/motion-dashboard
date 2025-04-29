package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"task-dashboard/internal/models"
)

// MotionClient handles interactions with the Motion API
type MotionClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewMotionClient creates a new API client for Motion
func NewMotionClient(apiKey string) *MotionClient {
	return &MotionClient{
		APIKey: apiKey,
		// Use exact URL format from JavaScript - no trailing slash
		BaseURL: "https://api.usemotion.com/v1",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchTasks retrieves all tasks from the Motion API
func (c *MotionClient) FetchTasks() ([]models.Task, error) {
	var allTasks []models.Task
	cursor := ""

	for {
		// Use exact URL format from JavaScript
		url := fmt.Sprintf("%s/tasks", c.BaseURL)
		if cursor != "" {
			url = fmt.Sprintf("%s?cursor=%s", url, cursor)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		// Set headers exactly as in JavaScript
		req.Header.Set("X-API-Key", c.APIKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Warn("Failed to close response body", "error", closeErr)
		}

		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
		}

		var response models.TasksResponse
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return nil, fmt.Errorf("error unmarshalling response: %w", err)
		}

		allTasks = append(allTasks, response.Tasks...)

		// Check if there are more pages
		if response.Meta.NextCursor == "" {
			break
		}
		cursor = response.Meta.NextCursor
	}

	return allTasks, nil
}

// ValidateAPIKey attempts to make a simple API call to verify the API key
func (c *MotionClient) ValidateAPIKey() error {
	// Use exact URL format from JavaScript
	url := fmt.Sprintf("%s/tasks", c.BaseURL)

	slog.Info("Validating API key with URL", "url", url)
	if len(c.APIKey) > 4 {
		if os.Getenv("LOG_API_KEYS") == "true" {
			slog.Info("API Key (first 4 chars)", "first_chars", c.APIKey[:4])
		}
	} else {
		slog.Info("API Key is too short", "length", len(c.APIKey))
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers exactly as in JavaScript
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	if closeErr := resp.Body.Close(); closeErr != nil {
		slog.Warn("Failed to close response body", "error", closeErr)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API key validation failed: Status=%s, Body=%s", resp.Status, string(bodyBytes))
	}

	slog.Info("API key successfully validated!")
	return nil
}
