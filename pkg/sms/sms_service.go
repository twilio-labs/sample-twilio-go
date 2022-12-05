package sms

import (
	"fmt"

	"github.com/twilio-labs/sample-twilio-go/pkg/configuration"
	"github.com/twilio-labs/sample-twilio-go/pkg/message"
	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	twilioAPI "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

/*
 * Service for handling SMS communication
 */
type SMSService struct {
	client *twilio.RestClient
	logger *zap.Logger
	config *configuration.TwilioConfiguration
}

/*
 * Constructor
 */
func NewSMSService(client *twilio.RestClient, logger *zap.Logger, config *configuration.TwilioConfiguration) *SMSService {
	return &SMSService{client, logger, config}
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
	}
	return nil
}
