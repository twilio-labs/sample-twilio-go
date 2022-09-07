package controller

import (
	"net/http"
	"strings"
	"time"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/sms"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/voice"
	"github.com/gin-gonic/gin"
	client "github.com/twilio/twilio-go/client"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	COOKIE_TTL              = 14400 // 4 hours in seconds. The TTL for Twilio SMS cookies
	WAIT_SECONDS_UNTIL_CALL = 5
)

/*
 * Controller for review related request resources
 */
type ReviewController struct {
	smsSvc       *sms.SMSService
	voiceSvc     *voice.VoiceService
	reqValidator *client.RequestValidator
	baseURL      string
}

/*
 * Constructor
 */
func NewReviewController(smsSvc *sms.SMSService, voiceSvc *voice.VoiceService, reqValidator *client.RequestValidator, baseURL string) *ReviewController {
	return &ReviewController{smsSvc, voiceSvc, reqValidator, baseURL}
}

func (ctr *ReviewController) HandleSMS(c *gin.Context) {
	// Validate Request
	if !ctr.isValidRequest(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Parse request cookies for SMS conversation context
	greeted, err := c.Cookie("greeted")
	if err != nil {
		greeted = "false"
		c.SetCookie("greeted", "false", COOKIE_TTL, "/sms", "localhost", false, false)
	}
	participant, err := c.Cookie("participant")
	if err != nil {
		participant = ""
		c.SetCookie("participant", "", COOKIE_TTL, "/sms", "localhost", false, false)
	}

	// Parse request POST form body
	incomingPhoneNumber := c.PostForm("From")
	incomingBody := strings.TrimSpace(c.PostForm("Body"))

	if greeted == "false" {
		// Send greeting and participation invite
		ctr.smsSvc.SendGreeting(incomingPhoneNumber)
		ctr.smsSvc.SendInvite(incomingPhoneNumber)
		c.SetCookie("greeted", "true", COOKIE_TTL, "/sms", "localhost", false, false)
		c.AbortWithStatus(http.StatusOK)
		return
	}

	if participant == "false" {
		// Participation invite declined
		ctr.smsSvc.SendGoodbye(incomingPhoneNumber)
		resetContext(c)
		c.AbortWithStatus(http.StatusOK)
		return
	}

	if participant == "" {
		if strings.ToLower(incomingBody) == "yes" {
			// Participation invite accepted. Query for participant name.
			c.SetCookie("participant", "true", COOKIE_TTL, "/sms", "localhost", false, false)
			ctr.smsSvc.SendAcceptConfirmation(incomingPhoneNumber)
			ctr.smsSvc.SendAskForName(incomingPhoneNumber)
			c.AbortWithStatus(http.StatusOK)
			return
		} else if strings.ToLower(incomingBody) == "no" {
			// Participation invite declined.
			c.SetCookie("participant", "false", COOKIE_TTL, "/sms", "localhost", false, false)
			ctr.smsSvc.SendGoodbye(incomingPhoneNumber)
			resetContext(c)
			c.AbortWithStatus(http.StatusOK)
			return
		} else {
			// Unknown input. Send fallback message.
			ctr.smsSvc.SendInviteFallback(incomingPhoneNumber)
			c.AbortWithStatus(http.StatusOK)
			return
		}
	}

	if participant == "true" && len(incomingBody) > 0 {
		// Invite acccepted. Name query
		caser := cases.Title(language.AmericanEnglish)
		ctr.smsSvc.SendNamedGreeting(incomingPhoneNumber, caser.String(incomingBody))
		ctr.smsSvc.SendCallNotification(incomingPhoneNumber)
		time.Sleep(time.Duration(WAIT_SECONDS_UNTIL_CALL) * time.Second)
		ctr.voiceSvc.InitiateReviewCall(incomingPhoneNumber)
		resetContext(c)
		c.AbortWithStatus(http.StatusOK)
		return
	} else {
		// Send name query fallback message
		ctr.smsSvc.SendAskForNameFallback(incomingPhoneNumber)
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

func resetContext(c *gin.Context) {
	c.SetCookie("participant", "", COOKIE_TTL, "/sms", "localhost", false, false)
	c.SetCookie("greeted", "false", COOKIE_TTL, "/sms", "localhost", false, false)
	c.SetCookie("identity", "", COOKIE_TTL, "/sms", "localhost", false, false)
}
