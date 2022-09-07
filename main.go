package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/configuration"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/controller"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/sms"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/voice"
	"github.com/gin-gonic/gin"
	twilio "github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
)

const (
	ACCOUNT_SID_ENV        = "TWILIO_ACCOUNT_SID"
	ACCOUNT_AUTH_TOKEN_ENV = "TWILIO_AUTH_TOKEN"
)

func main() {
	// Process CLI and env args
	from := flag.String("from", "", "From phone number")
	baseUrl := flag.String("url", "", "Server Base URL")
	flag.Parse()

	if *from == "" {
		log.Fatal("Missing required from phone number arg.")
	}
	if *baseUrl == "" {
		log.Fatal("Missing required base URL arg.")
	}

	accountSID := os.Getenv(ACCOUNT_SID_ENV)
	authToken := os.Getenv(ACCOUNT_AUTH_TOKEN_ENV)

	// Initialize Twilio REST client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	// Initialize application service(s)
	config := &configuration.TwilioConfiguration{
		AccountSID:           accountSID,
		AccountPhoneNumber:   *from,
		BaseURL:              *baseUrl,
		StatusCallbackPath:   "/call-event",
		StatusCallbackMethod: "POST",
	}
	smsSvc := sms.NewSMSService(client, config)
	voiceSvc := voice.NewVoiceService(client, config)
	reqValidator := twilioClient.NewRequestValidator(authToken)

	// Initialize application controller(s)
	reviewCtr := controller.NewReviewController(smsSvc, voiceSvc, &reqValidator, config.BaseURL)

	// Initialize Gin request routing and start Gin web server
	r := gin.Default()
	r.POST("/sms", reviewCtr.HandleSMS)
	r.POST("/call-event", reviewCtr.HandleCallEvent)
	r.GET("/call-total", func(c *gin.Context) {
		logs, err := voiceSvc.RetrieveCallLogs()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "Total Calls: %d", len(logs))
	})
	r.Run()
}
