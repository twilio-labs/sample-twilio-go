package message

import (
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestGetReviewGreetingAndInstructionsTwiML(t *testing.T) {
	// Act
	result, error := GetReviewGreetingAndInstructionsTwiML()
	if error != nil {
		t.Error(error)
	}

	// Assert that result contains the expected TwiML
	assert.Equal(t, true, strings.Contains(result, "Ahoy! Greetings from Twilio Resorts and Spas"))

}
