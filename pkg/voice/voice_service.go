package voice

import (
	"fmt"

	"github.com/twilio-labs/sample-twilio-go/pkg/configuration"
	"github.com/twilio-labs/sample-twilio-go/pkg/message"
	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	twilioAPI "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
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
	logger *zap.Logger
	config *configuration.TwilioConfiguration
}

/*
 * Constructor
 */
func NewVoiceService(client *twilio.RestClient, logger *zap.Logger, config *configuration.TwilioConfiguration) *VoiceService {
	return &VoiceService{client, logger, config}
}

func (svc *VoiceService) InitiateReviewCall(to string) error {
	params := &twilioAPI.CreateCallParams{}
	params.SetTo(to)
	params.SetFrom(svc.config.AccountPhoneNumber)
	params.SetStatusCallback(svc.config.BaseURL + svc.config.StatusCallbackPath)
	params.SetStatusCallbackEvent([]string{"completed"})
	params.SetStatusCallbackMethod(svc.config.StatusCallbackMethod)

	twiml, err := message.GetReviewGreetingAndInstructionsTwiML()
	if err != nil {
		return fmt.Errorf("[InitiateReviewCall] Failed to generate call TwiML. Error: %w", err)
	}
	params.SetTwiml(twiml)

	svc.logCallParams(params)

	_, err = svc.client.Api.CreateCall(params)
	if err != nil {
		svc.logTwilioError("Failed to create call", err)
		return fmt.Errorf("[InitiateReviewCall] Failed to create review call. Error: %w", err)
	}
	return nil
}

func (svc *VoiceService) RetrieveCallLogs() ([]*ReviewCallRecord, error) {
	var logs []*ReviewCallRecord
	params := &twilioAPI.ListCallParams{}
	records, err := svc.client.Api.ListCall(params)
	if err != nil {
		svc.logTwilioError("Failed to retrieve call logs", err)
		return nil, fmt.Errorf("[RetrieveCallLogs] Failed to retrieve call logs. Error: %w", err)
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

func (svc *VoiceService) logTwilioError(msg string, err error) {
	twilioError := err.(*twilioClient.TwilioRestError)
	svc.logger.Error(msg,
		zap.Int("code", twilioError.Code),
		zap.String("message", twilioError.Message),
		zap.Int("status", twilioError.Status),
		zap.String("moreInfo", twilioError.MoreInfo),
		zap.Any("details", twilioError.Details))
}

func (svc *VoiceService) logCallParams(params *twilioAPI.CreateCallParams) {
	svc.logger.Debug("Voice call parameters",
		zap.String("To", *params.To),
		zap.String("From", *params.From),
		zap.String("StatusCallback", *params.StatusCallback),
		zap.Strings("StatusCallbackEvent", *params.StatusCallbackEvent),
		zap.String("StatusCallbackMethod", *params.StatusCallbackMethod),
		zap.String("Twiml", *params.Twiml))
}
