package queries

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Client struct {
	endpoint string
	client   *http.Client
}

// New creates a new GraphQL client
func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Execute sends a GraphQL query and unmarshals the result into the generic type T
func Execute[T any](ctx context.Context, c *Client, query string, variables map[string]interface{}) (*T, error) {
	// Construct payload
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Perform request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 status code: %d\nresponse body: %s", resp.StatusCode, string(respBody))
	}

	// Response wrapper
	type responseWrapper struct {
		Data   T             `json:"data"`
		Errors []interface{} `json:"errors"` // Optional: could be enhanced
	}

	var wrapper responseWrapper
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(wrapper.Errors) > 0 {
		return nil, fmt.Errorf("graphql errors: %+v", wrapper.Errors)
	}

	return &wrapper.Data, nil
}

func QueryYearly[T any](ctx context.Context, client *Client, query string, vars map[string]interface{}, start, end time.Time) ([]*T, error) {
	var queries []*T

	for curr := start; curr.Before(end); {
		// Define the next interval (1 year ahead)
		next := curr.AddDate(1, 0, 0)

		if next.After(end) {
			next = end
		}

		vars["from"] = curr.Format(time.RFC3339)
		vars["to"] = next.Format(time.RFC3339)

		query, err := Execute[T](ctx, client, query, vars)
		if err != nil {
			return nil, err
		}
		queries = append(queries, query)

		curr = next
	}

	vars["from"] = start
	vars["to"] = time.Now()
	return queries, nil
}

func QueryRecursively[T any](ctx context.Context, client *Client, query string, vars map[string]interface{}, extract PageInfoExtractor[T]) ([]*T, error) {
	var queries []*T

	for {
		query, err := Execute[T](ctx, client, query, vars)
		if err != nil {
			log.Fatal(err)
		}
		queries = append(queries, query)

		pageInfo := extract(query)
		if !pageInfo.HasNextPage {
			break
		}
		vars["cursor"] = pageInfo.EndCursor
	}

	vars["cursor"] = nil
	return queries, nil
}
