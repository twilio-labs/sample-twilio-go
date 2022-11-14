package controller

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twilio-labs/sample-twilio-go/pkg/db"
	"github.com/twilio-labs/sample-twilio-go/pkg/sms"
	"github.com/twilio-labs/sample-twilio-go/pkg/voice"
	client "github.com/twilio/twilio-go/client"
)

const (
	COOKIE_TTL              = 14400 // 4 hours in seconds. The TTL for Twilio SMS cookies
	WAIT_SECONDS_UNTIL_CALL = 8
)

/*
 * Controller for review related request resources
 */
type ReviewController struct {
	ctx          context.Context
	db           *db.DB
	smsSvc       *sms.SMSService
	voiceSvc     *voice.VoiceService
	reqValidator *client.RequestValidator
	baseURL      string
}

/*
 * Constructor
 */
func NewReviewController(ctx context.Context, db *db.DB, smsSvc *sms.SMSService, voiceSvc *voice.VoiceService, reqValidator *client.RequestValidator, baseURL string) *ReviewController {
	return &ReviewController{ctx, db, smsSvc, voiceSvc, reqValidator, baseURL}
}

func (ctr *ReviewController) StartReviewCampaign(c *gin.Context) {
	customers, err := ctr.db.GetCustomers(ctr.ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(customers))
	for i := 0; i < len(customers); i++ {
		go func(i int) {
			defer wg.Done()
			c := customers[i]
			// This should probably handle the error in a multi-threaded way too
			_ = ctr.sendGreetingAndInvite(c.FirstName, c.PhoneNumber)
		}(i)
	}
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Success",
		"invites": len(customers),
	})
}

func (ctr *ReviewController) HandleSMS(c *gin.Context) {
	// Validate Request
	if !ctr.isValidRequest(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Parse request POST form body
	incomingPhoneNumber := c.PostForm("From")
	incomingBody := strings.TrimSpace(c.PostForm("Body"))

	// Should probably check if the "From" number is a registered customer

	if strings.ToLower(incomingBody) == "yes" {
		// Participation invite accepted. Query for participant name.
		err := ctr.sendInviteComfirmation(incomingPhoneNumber)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		err = ctr.sendCallNotificationAndInitiateCall(incomingPhoneNumber)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatus(http.StatusOK)
		return
	} else if strings.ToLower(incomingBody) == "no" {
		// Participation invite declined.
		err := ctr.smsSvc.SendGoodbye(incomingPhoneNumber)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatus(http.StatusOK)
		return
	} else {
		// Unknown input. Send fallback message.
		err := ctr.smsSvc.SendInviteFallback(incomingPhoneNumber)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatus(http.StatusOK)
		return
	}
}

func (ctr *ReviewController) HandleCallEvent(c *gin.Context) {
	// Validate Request
	if !ctr.isValidRequest(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	receiverPhone := c.PostForm("To")
	ctr.smsSvc.SendThankYou(receiverPhone)
	c.AbortWithStatus(http.StatusOK)
}

func (ctr *ReviewController) isValidRequest(c *gin.Context) bool {
	// Validate request with Twilio SDK request validator
	url := ctr.baseURL + c.Request.RequestURI
	signatureHeader := c.Request.Header["X-Twilio-Signature"]
	c.Request.ParseForm()
	params := make(map[string]string)
	for k, v := range c.Request.PostForm {
		params[k] = v[0]
	}
	if len(signatureHeader) > 0 {
		return ctr.reqValidator.Validate(url, params, signatureHeader[0])
	}
	return false
}

func (ctr *ReviewController) sendGreetingAndInvite(name, to string) error {
	if err := ctr.smsSvc.SendGreeting(name, to); err != nil {
		return err
	}
	if err := ctr.smsSvc.SendInvite(to); err != nil {
		return err
	}
	return nil
}

func (ctr *ReviewController) sendInviteComfirmation(to string) error {
	if err := ctr.smsSvc.SendAcceptConfirmation(to); err != nil {
		return err
	}
	return nil
}

func (ctr *ReviewController) sendCallNotificationAndInitiateCall(to string) error {
	if err := ctr.smsSvc.SendCallNotification(to); err != nil {
		return err
	}
	time.Sleep(time.Duration(WAIT_SECONDS_UNTIL_CALL) * time.Second)
	if err := ctr.voiceSvc.InitiateReviewCall(to); err != nil {
		return err
	}
	return nil
}
