package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/sms"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/voice"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	twilioClient "github.com/twilio/twilio-go/client"
)

func TestNewReviewController(t *testing.T) {
	// Arrange
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(sms, voice, clientRequestValidator, "http://localhost")

	// Assert
	assert.NotNil(t, ctr)
}

func GetTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func TestReviewController_HandleSMS(t *testing.T) {
	// Arrange
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(sms, voice, clientRequestValidator, "http://localhost")

	ctx := GetTestGinContext()

	// Act
	ctr.HandleSMS(ctx)

	// Assert
	assert.Equal(t, 403, ctx.Writer.Status())
}

func TestReviewController_HandleCallHandler(t *testing.T) {
	// Arrange
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(sms, voice, clientRequestValidator, "http://localhost")

	ctx := GetTestGinContext()

	// Act
	ctr.HandleCallEvent(ctx)

	// Assert
	assert.Equal(t, 403, ctx.Writer.Status())
}

// check isValidRequest
func TestReviewController_isValidRequest(t *testing.T) {
	// Arrange
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(sms, voice, clientRequestValidator, "http://localhost")

	ctx := GetTestGinContext()

	// Act
	ctr.isValidRequest(ctx)

	// Assert
	assert.Equal(t, 200, ctx.Writer.Status())
}

// test resetContext
func TestReviewController_resetContext(t *testing.T) {
	ctx := GetTestGinContext()

	// Act
	resetContext(ctx)

	// Assert
	assert.Equal(t, 200, ctx.Writer.Status())
}
