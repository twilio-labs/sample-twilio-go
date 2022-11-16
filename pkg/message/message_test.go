package message

import (
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

// test GetHelloMessage
func TestGetHelloMessage(t *testing.T) {
	// Arrange
	name := "John"

	// Act
	result := GetHelloMessage(name)

	// Assert
	assert.Equal(t, "Hello, John!", result)
}

func TestGetReviewGreetingAndInstructionsTwiML(t *testing.T) {
	// Act
	result, error := GetReviewGreetingAndInstructionsTwiML()
	if error != nil {
		t.Error(error)
	}

	// Assert that result contains the expected TwiML
	assert.Equal(t, true, strings.Contains(result, "Ahoy! Greetings from Twilio Resorts and Spas"))

}
