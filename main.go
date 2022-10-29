package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/configuration"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/controller"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/sms"
	"code.hq.twilio.com/twilio/review-rewards-example-app/pkg/voice"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	twilio "github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ACCOUNT_SID_ENV        = "TWILIO_ACCOUNT_SID"
	ACCOUNT_AUTH_TOKEN_ENV = "TWILIO_AUTH_TOKEN"
	PHONE_NUMBER_ENV       = "TWILIO_PHONE_NUMBER"
	BASE_URL_ENV           = "BASE_URL"
	LOG_LEVEL_ENV          = "LOG_LEVEL"
)

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "api",
		Name:       "latency_seconds",
		Help:       "Latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)

func prometheusGinMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	latency.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(time.Since(start).Seconds())
}

func init() {
	prometheus.MustRegister(latency)
}

func main() {
	// Process CLI and env args
	from := os.Getenv(PHONE_NUMBER_ENV)
	baseUrl := os.Getenv(BASE_URL_ENV)
	logLevel := os.Getenv(LOG_LEVEL_ENV)
	if from == "" {
		log.Fatal("Missing required from phone number arg.")
	}
	if baseUrl == "" {
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
		AccountPhoneNumber:   from,
		BaseURL:              baseUrl,
		StatusCallbackPath:   "/call-event",
		StatusCallbackMethod: "POST",
	}
	logger, err := initializeLogger(logLevel)
	if err != nil {
		log.Fatal("Failed to initialize logger. Error: ", err)
	}
	defer logger.Sync() // flushes buffer, if any
	smsSvc := sms.NewSMSService(client, logger, config, latency)
	voiceSvc := voice.NewVoiceService(client, logger, config)
	reqValidator := twilioClient.NewRequestValidator(authToken)

	// Initialize application controller(s)
	reviewCtr := controller.NewReviewController(smsSvc, voiceSvc, &reqValidator, config.BaseURL)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(prometheusGinMiddleware)

	// Request routing
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.POST("/sms", reviewCtr.HandleSMS)
	r.POST("/call-event", reviewCtr.HandleCallEvent)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
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
