package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twilio-labs/sample-twilio-go/pkg/configuration"
	"github.com/twilio-labs/sample-twilio-go/pkg/controller"
	"github.com/twilio-labs/sample-twilio-go/pkg/db"
	"github.com/twilio-labs/sample-twilio-go/pkg/metric"
	"github.com/twilio-labs/sample-twilio-go/pkg/sms"
	"github.com/twilio-labs/sample-twilio-go/pkg/voice"
	twilio "github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ACCOUNT_SID_ENV        = "TWILIO_ACCOUNT_SID"
	ACCOUNT_AUTH_TOKEN_ENV = "TWILIO_AUTH_TOKEN"
	PHONE_NUMBER_ENV       = "TWILIO_PHONE_NUMBER"
	BASE_URL_ENV           = "BASE_URL"
	LOG_LEVEL_ENV          = "LOG_LEVEL"
	TEMPLATE_DIR_ENV       = "TEMPLATE_DIR"
	ASSET_DIR_ENV          = "ASSET_DIR"
)

var (
	// Defines the quantile rank estimates with their respective
	// absolute error. If Objectives[q] = e, then the value reported for q
	// will be the φ-quantile value for some φ between q-e and q+e.
	defaultLatencyObjectives = map[float64]float64{
		0.5:  0.05,
		0.9:  0.01,
		0.99: 0.001,
	}
)

// Prometheus metrics
var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "api",
		Name:       "latency_seconds",
		Help:       "Latency distributions.",
		Objectives: defaultLatencyObjectives,
	},
	[]string{"method", "path"},
)
var invitesSentCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "api",
		Name:      "invites_sent",
		Help:      "Review invites sent.",
	},
)
var invitesAcceptedCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "api",
		Name:      "invites_accepted",
		Help:      "Review invites accepted.",
	},
)
var invitesDeclinedCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "api",
		Name:      "invites_declined",
		Help:      "Review invites declined.",
	},
)

var tracer = otel.GetTracerProvider().Tracer("twilio-go-at-scale")

func prometheusGinMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	latency.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(time.Since(start).Seconds())

	// tracer
	ctx, span := tracer.Start(c.Request.Context(), "prometheusGinMiddleware")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)
}

func init() {
	prometheus.MustRegister(latency)
	prometheus.MustRegister(invitesSentCounter)
	prometheus.MustRegister(invitesAcceptedCounter)
	prometheus.MustRegister(invitesDeclinedCounter)
}

func main() {
	// Setup Tracer
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Process CLI and env args
	templateDir := os.Getenv(TEMPLATE_DIR_ENV)
	assetDir := os.Getenv(ASSET_DIR_ENV)
	from := os.Getenv(PHONE_NUMBER_ENV)
	baseUrl := os.Getenv(BASE_URL_ENV)
	logLevel := os.Getenv(LOG_LEVEL_ENV)
	if len(templateDir) == 0 {
		templateDir = "../../pkg/template"
	}
	if len(assetDir) == 0 {
		assetDir = "../../asset"
	}
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
	db, err := db.InitializeDB()
	if err != nil {
		log.Fatal("Failed to initialize the database connection. Error: ", err)
	}

	// Initialize Metrics
	inviteMetrics := metric.NewInviteMetrics(invitesSentCounter,
		invitesAcceptedCounter,
		invitesDeclinedCounter)

	smsSvc := sms.NewSMSService(client, logger, config, latency)
	voiceSvc := voice.NewVoiceService(client, logger, config)
	reqValidator := twilioClient.NewRequestValidator(authToken)

	// Initialize application controller(s)
	reviewCtr := controller.NewReviewController(ctx, db, smsSvc, voiceSvc, &reqValidator, config.BaseURL, inviteMetrics)
	registerCtr := controller.NewRegisterController(ctx, db)
	controlPanelCtr := controller.NewControlPanelController()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(prometheusGinMiddleware)

	// Initialize HTML templates, static assets
	initTemplates(r, templateDir)
	initAssets(r, assetDir)

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
	r.GET("/register", registerCtr.GET)
	r.POST("/register", registerCtr.POST)
	r.GET("/campaigns-control-panel", controlPanelCtr.GET)
	r.POST("/campaign-start", reviewCtr.StartReviewCampaign)
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.Run()
}

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// export traces to Zap logger
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// Identify your application using resource detection
	res, err := resource.New(ctx,
		// Use the GCP resource detector to detect information about the GCP platform
		//resource.WithDetectors(gcp.NewDetector()),
		// Keep the default detectors
		resource.WithTelemetrySDK(),
		// Add your own custom attributes to identify your application
		resource.WithAttributes(
			semconv.ServiceNameKey.String("my-application"),
		),
	)
	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	defer tp.ForceFlush(ctx) // flushes any pending spans
	otel.SetTracerProvider(tp)
	return tp, nil
}

func initTemplates(r *gin.Engine, templateDir string) {
	r.LoadHTMLGlob(templateDir + "/*.html")
}

func initAssets(r *gin.Engine, assetDir string) {
	r.Static("/asset", assetDir)
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
