package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/twilio-labs/sample-twilio-go/pkg/db"
	"github.com/twilio-labs/sample-twilio-go/pkg/metric"
	"github.com/twilio-labs/sample-twilio-go/pkg/sms"
	"github.com/twilio-labs/sample-twilio-go/pkg/voice"
	twilioClient "github.com/twilio/twilio-go/client"
)

func TestNewReviewController(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db, _ := db.InitializeDB()
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}
	metrics := &metric.InviteMetrics{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(ctx, db, sms, voice, clientRequestValidator, "http://localhost", metrics)

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
	ctx := context.Background()
	db, _ := db.InitializeDB()
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}
	metrics := &metric.InviteMetrics{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(ctx, db, sms, voice, clientRequestValidator, "http://localhost", metrics)

	testGinCtx := GetTestGinContext()

	// Act
	ctr.HandleSMS(testGinCtx)

	// Assert
	assert.Equal(t, 403, testGinCtx.Writer.Status())
}

func TestReviewController_HandleCallHandler(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db, _ := db.InitializeDB()
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}
	metrics := &metric.InviteMetrics{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(ctx, db, sms, voice, clientRequestValidator, "http://localhost", metrics)

	testCtx := GetTestGinContext()

	// Act
	ctr.HandleCallEvent(testCtx)

	// Assert
	assert.Equal(t, 403, testCtx.Writer.Status())
}

// check isValidRequest
func TestReviewController_isValidRequest(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db, _ := db.InitializeDB()
	sms := &sms.SMSService{}
	voice := &voice.VoiceService{}
	metrics := &metric.InviteMetrics{}

	clientRequestValidator := &twilioClient.RequestValidator{}
	ctr := NewReviewController(ctx, db, sms, voice, clientRequestValidator, "http://localhost", metrics)

	testCtx := GetTestGinContext()

	// Act
	ctr.isValidRequest(testCtx)

	// Assert
	assert.Equal(t, 200, testCtx.Writer.Status())
}
