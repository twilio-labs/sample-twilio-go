package cdp

import (
	"context"

	analytics "github.com/segmentio/analytics-go/v3"
)

// utlize the segment client to send events to the customer data platform
type Client struct {
	client analytics.Client
}

// NewClient creates a new client for the customer data platform
func NewClient(writeKey string) Client {
	c := analytics.New(writeKey)
	return Client{client: c}
}

// Identify sends an identify event to the customer data platform with the given userid and properties
func (c *Client) Identify(ctx context.Context, userID string, properties map[string]interface{}) error {
	return c.client.Enqueue(analytics.Identify{
		UserId:     userID,
		Traits: properties,
	})
}

// Track sends a track event to the customer data platform with the given event name, userid and properties
func (c *Client) Track(ctx context.Context, event string, userID string, properties map[string]interface{}) error {
	return c.client.Enqueue(analytics.Track{
		Event:      event,
		UserId:     userID,
		Properties: properties,
	})
}

// Close closes the client connection to the customer data platform
func (c *Client) Close() {
	c.client.Close()
}