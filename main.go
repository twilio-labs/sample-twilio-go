package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/twilio-labs/sample-twilio-go/pkg/configuration"
	"github.com/twilio-labs/sample-twilio-go/pkg/controller"
	"github.com/twilio-labs/sample-twilio-go/pkg/sms"
	"github.com/twilio-labs/sample-twilio-go/pkg/voice"
	twilio "github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ACCOUNT_SID_ENV        = "TWILIO_ACCOUNT_SID"
	ACCOUNT_AUTH_TOKEN_ENV = "TWILIO_AUTH_TOKEN"
)

func main() {
	// Process CLI and env args
	from := flag.String("from", "", "From phone number")
	baseUrl := flag.String("url", "", "Server Base URL")
	logLevel := flag.String("loglevel", "info", "Log level")
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
	logger, err := initializeLogger(*logLevel)
	if err != nil {
		log.Fatal("Failed to initialize logger. Error: ", err)
	}
	defer logger.Sync() // flushes buffer, if any
	smsSvc := sms.NewSMSService(client, logger, config)
	voiceSvc := voice.NewVoiceService(client, logger, config)
	reqValidator := twilioClient.NewRequestValidator(authToken)

	// Initialize application controller(s)
	reviewCtr := controller.NewReviewController(smsSvc, voiceSvc, &reqValidator, config.BaseURL)

	r := gin.Default()

	// Request routing
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
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.Run()
}

func initializeLogger(logLevel string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	switch logLevel {
	case "debug":
		cfg.Level.SetLevel(zapcore.DebugLevel)
	case "info":
		cfg.Level.SetLevel(zapcore.InfoLevel)
	case "warn":
		cfg.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		cfg.Level.SetLevel(zapcore.ErrorLevel)
	case "panic":
		cfg.Level.SetLevel(zapcore.PanicLevel)
	case "fatal":
		cfg.Level.SetLevel(zapcore.FatalLevel)
	default:
		return nil, fmt.Errorf("invalid log level argument")
	}
	return cfg.Build()
}
