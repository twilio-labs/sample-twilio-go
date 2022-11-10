package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/twilio-labs/sample-twilio-go/pkg/configuration"
	"github.com/twilio-labs/sample-twilio-go/pkg/message"
	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	twilioAPI "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

// Path: pkg/sms/sms_service.go
type PublishContent struct {
	ToNumber   string `json:"to"`
	Message    string `json:"message"`
}

/*
 * Service for handling SMS communication
*/
type SMSService struct {
	client  *twilio.RestClient
	logger  *zap.Logger
	config  *configuration.TwilioConfiguration
	latency *prometheus.SummaryVec
}

/*
 * Constructor
 */
func NewSMSService(client *twilio.RestClient, logger *zap.Logger, config *configuration.TwilioConfiguration, latency *prometheus.SummaryVec) *SMSService {
	return &SMSService{client, logger, config, latency}
}

func (svc *SMSService) SendGreeting(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.GREETING,
		"[SendGreeting] Failed to send greeting")
}

func (svc *SMSService) SendInvite(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.PARTICIPATION_INVITE,
		"[SendInvite] Failed to send review invite")
}

func (svc *SMSService) SendAcceptConfirmation(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.PARTICIPATION_ACCEPT_RESPONSE,
		"[SendAcceptConfirmation] Failed to send invite accept confirmation")
}

func (svc *SMSService) SendInviteFallback(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.PARTICIPATION_INVITE_FALLBACK,
		"[SendInviteFallback] Failed to send invite fallback")
}

func (svc *SMSService) SendAskForName(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.ASK_FOR_NAME,
		"[SendAskForName] Failed to send name query")
}

func (svc *SMSService) SendAskForNameFallback(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.ASK_FOR_NAME,
		"[SendAskForNameFallback] Failed to send name query fallback")
}

func (svc *SMSService) SendNamedGreeting(to, name string) error {
	body := message.GetHelloMessage(name)
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		body,
		"[SendNamedGreeting] Failed to send named greeting")
}

func (svc *SMSService) SendCallNotification(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.CALL_NOTIFICATION,
		"[SendCallNotification] Failed to send call notification")
}

func (svc *SMSService) SendThankYou(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.PARTICIPATION_THANKYOU,
		"[SendThankYou] Failed to send thank you")
}

func (svc *SMSService) SendGoodbye(to string) error {
	return svc.sendMessage(to,
		svc.config.AccountPhoneNumber,
		message.GOODBYE,
		"[SendGoodbye] Failed to send goodbye")
}

func (svc *SMSService) sendMessage(to, from, body, errMsg string) error {
	params := &twilioAPI.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(body)

	// measure the latency of the request
	start := time.Now()

	// Logging the request parameters sent to the SDK client
	svc.logger.Debug("SMS message parameters",
		zap.String("To", *params.To),
		zap.String("From", *params.From),
		zap.String("Body", *params.Body))

	_, err := svc.client.Api.CreateMessage(params)

	// Debug logging errors returned from Twilio Messaging
	if err != nil {
		twilioError := err.(*twilioClient.TwilioRestError)
		svc.logger.Error("Failed to send SMS message",
			zap.Int("code", twilioError.Code),
			zap.String("message", twilioError.Message),
			zap.Int("status", twilioError.Status),
			zap.String("moreInfo", twilioError.MoreInfo),
			zap.Any("details", twilioError.Details))

		return fmt.Errorf("%s. Error: %w", errMsg, err)
	} else {
		// Logging the response from Twilio Messaging
		svc.logger.Debug("SMS message sent successfully")
		latency := time.Since(start).Seconds()
		svc.logger.Debug("Twilio SMS message latency", zap.Float64("latency", latency))
		svc.latency.WithLabelValues("TWILIO", "SMS").Observe(latency)
				// Send to Pub/Sub
				var wg sync.WaitGroup
				projectId := os.Getenv("PROJECT_ID")
				topic := os.Getenv("TOPIC")
		
				client, err := pubsub.NewClient(context.Background(), projectId)
				if err != nil {
					log.Fatalf("Failed to create client: %v", err)
				}
				defer client.Close()
		
				t := client.Topic(topic)
				// setup json with message content
				msg := PublishContent{
					ToNumber: to,
					Message: body,
				}
		
				msgJson, _ := json.Marshal(msg)
				result := t.Publish(context.Background(), &pubsub.Message{
					Data: msgJson,
				})
		
				wg.Add(1)
		
				go func(res *pubsub.PublishResult) {
					defer wg.Done()
		
					_, err := res.Get(context.Background())
					if err != nil {
						fmt.Fprintf(os.Stdout, "Failed to publish: %v", err)
						return
					}
				}(result)
				
				wg.Wait()
	}
	return nil
}
