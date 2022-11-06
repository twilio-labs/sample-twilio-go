package cdp

// test the client.go file

import (
	"context"
	"testing"

	"github.com/bmizerany/assert"
)

// test the NewClient function
func TestNewClient(t *testing.T) {
	// Arrange
	writeKey := "test"

	// Act
	result := NewClient(writeKey)

	// Assert
	assert.Equal(t, true, result != Client{})
}

// test the Identify function
func TestIdentify(t *testing.T) {

	// Arrange
	writeKey := "test"
	c := NewClient(writeKey)
	userID := "test"
	properties := map[string]interface{}{"test": "test"}

	// Act
	result := c.Identify(context.Background(), userID, properties)

	// Assert
	assert.Equal(t, nil, result)
}

// test the Track function
func TestTrack(t *testing.T) {
	// Arrange
	writeKey := "test"
	c := NewClient(writeKey)
	event := "test"
	userID := "test"
	properties := map[string]interface{}{"test": "test"}

	// Act
	result := c.Track(context.Background(), event, userID, properties)

	// Assert
	assert.Equal(t, nil, result)
}


