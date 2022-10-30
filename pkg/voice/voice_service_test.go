package voice

import (
	"testing"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/configuration"
	"github.com/twilio/twilio-go"
	"go.uber.org/zap"
)

// test newVoiceService
func TestNewVoiceService(t *testing.T) {
	// Arrange
	twilioClient := &twilio.RestClient{}
	twilioConfiguration := &configuration.TwilioConfiguration{}

	// setup zap logger noop
	logger := zap.NewNop()

	voiceService := NewVoiceService(twilioClient, logger, twilioConfiguration)

	if voiceService == nil {
		t.Errorf("Expected voiceService to not be nil")
	}
}
