package voice

import (
	"fmt"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/configuration"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/message"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type ReviewCallRecord struct {
	SID         string
	To          string
	DateCreated string
	Duration    string
}

/*
 * Service for handling voice communication
 */
type VoiceService struct {
	client *twilio.RestClient
	config *configuration.TwilioConfiguration
}

/*
 * Constructor
 */
func NewVoiceService(client *twilio.RestClient, config *configuration.TwilioConfiguration) *VoiceService {
	return &VoiceService{client, config}
}

func (svc *VoiceService) InitiateReviewCall(to string) error {
	params := &openapi.CreateCallParams{}
	params.SetTo(to)
	params.SetFrom(svc.config.AccountPhoneNumber)
	params.SetStatusCallback(svc.config.BaseURL + svc.config.StatusCallbackPath)
	params.SetStatusCallbackEvent([]string{"completed"})
	params.SetStatusCallbackMethod(svc.config.StatusCallbackMethod)
	params.SetTwiml(message.GetReviewGreetingAndInstructionsTwiML())
	_, err := svc.client.Api.CreateCall(params)
	if err != nil {
		return fmt.Errorf("[InitiateReviewCall] Failed to create review call.\n Error: %w", err)
	}
	return nil
}

func (svc *VoiceService) RetrieveCallLogs() ([]*ReviewCallRecord, error) {
	var logs []*ReviewCallRecord
	params := &openapi.ListCallParams{}
	records, err := svc.client.Api.ListCall(params)
	if err != nil {
		return nil, fmt.Errorf("[RetrieveCallLogs] Failed to retrieve call logs.\n Error: %w", err)
	}
	for _, v := range records {
		log := ReviewCallRecord{
			SID:         *v.Sid,
			To:          *v.To,
			DateCreated: *v.DateCreated,
			Duration:    *v.Duration,
		}
		logs = append(logs, &log)
	}
	return logs, nil
}
