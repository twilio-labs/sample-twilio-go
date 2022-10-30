package sms

import (
	"testing"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/configuration"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/twilio/twilio-go"
	"go.uber.org/zap"
)

func TestNewSMSService(t *testing.T) {
	// Arrange
	twilioClient := &twilio.RestClient{}
	twilioConfiguration := &configuration.TwilioConfiguration{}

	// setup zap logger noop
	logger := zap.NewNop()
	
	// setup latency prometheus.SummaryVec
	latency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "sms_latency",
			Help: "Latency of SMS requests",
		},
		[]string{"method"},
	)

	smsService := NewSMSService(twilioClient, logger, twilioConfiguration, latency)

	if smsService == nil {
		t.Errorf("Expected smsService to not be nil")
	}
}