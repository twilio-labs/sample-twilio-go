package sms

import (
	"fmt"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/configuration"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/message"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

/*
 * Service for handling SMS communication
 */
type SMSService struct {
	client *twilio.RestClient
	config *configuration.TwilioConfiguration
}

/*
 * Constructor
 */
func NewSMSService(client *twilio.RestClient, config *configuration.TwilioConfiguration) *SMSService {
	return &SMSService{client, config}
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
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(body)
	_, err := svc.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("%s.\nError: %w", errMsg, err)
	}
	return nil
}
