package configuration

import (
	"testing"

	"github.com/bmizerany/assert"
)

// test twilio_configuration.go
func TestGenerateUserID(t *testing.T) {
	// Act
	result := GenerateUserID()

	// Assert
	assert.Equal(t, true, result != "")
}
