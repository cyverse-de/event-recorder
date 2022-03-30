package main

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/DavidGamba/go-getoptions"
	"github.com/cyverse-de/configurate"
	"github.com/cyverse-de/event-recorder/common"
	"github.com/cyverse-de/event-recorder/db"
	"github.com/cyverse-de/event-recorder/handlers"
	"github.com/cyverse-de/event-recorder/handlerset"
	"github.com/cyverse-de/event-recorder/logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var log = logging.Log.WithFields(logrus.Fields{"package": "main"})

// commandLineOptionValues represents the values of the command-line options that were passed on the command line when
// this service was invoked.
type commandLineOptionValues struct {
	Config string
}

// parseCommandLine parses the command line and returns an options structure containing command-line options and
// parameters.
func parseCommandLine() *commandLineOptionValues {
	optionValues := &commandLineOptionValues{}
	opt := getoptions.New()

	// Default option values.
	defaultConfigPath := "/etc/iplant/de/jobservices.yml"

	// Define the command-line options.
	opt.Bool("help", false, opt.Alias("h", "?"))
	opt.StringVar(&optionValues.Config, "config", defaultConfigPath,
		opt.Alias("c"),
		opt.Description("the path to the configuration file"))

	// Parse the command line, handling requests for help and usage errors.
	_, err := opt.Parse(os.Args[1:])
	if opt.Called("help") {
		fmt.Fprintf(os.Stderr, opt.Help())
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		fmt.Fprintf(os.Stderr, opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}

	return optionValues
}

func jaegerTracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("event-recorder"),
		)),
	)

	return tp, nil
}

func main() {
	var tracerProvider *tracesdk.TracerProvider

	otelTracesExporter := os.Getenv("OTEL_TRACES_EXPORTER")
	if otelTracesExporter == "jaeger" {
		jaegerEndpoint := os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT")
		if jaegerEndpoint == "" {
			log.Warn("Jaeger set as OpenTelemetry trace exporter, but no Jaeger endpoint configured.")
		} else {
			tp, err := jaegerTracerProvider(jaegerEndpoint)
			if err != nil {
				log.Fatal(err)
			}
			tracerProvider = tp
			otel.SetTracerProvider(tp)
			otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		}
	}

	if tracerProvider != nil {
		tracerCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		defer func(tracerContext context.Context) {
			ctx, cancel := context.WithTimeout(tracerContext, time.Second*5)
			defer cancel()
			if err := tracerProvider.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}(tracerCtx)
	}

	// Parse the command-line.
	optionValues := parseCommandLine()

	// Read in the configuration file.
	cfg, err := configurate.InitDefaults(optionValues.Config, configurate.JobServicesDefaults)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the AMQP settings.
	amqpSettings := &common.AMQPSettings{
		URI:          cfg.GetString("amqp.uri"),
		ExchangeName: cfg.GetString("amqp.exchange.name"),
		ExchangeType: cfg.GetString("amqp.exchange.type"),
	}

	// Initialize the database connection.
	databaseURI := cfg.GetString("notifications.db.uri")
	db, err := db.InitDatabase("postgres", databaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the email address to use for support requests.
	supportEmail := cfg.GetString("email.request")

	// Initialize the message handlers.
	messageHandlers, err := handlers.InitMessageHandlers(db, amqpSettings)
	if err != nil {
		log.Fatal(err)
	}

	// Create the message handler set.
	handlerSet, err := handlerset.New(amqpSettings, supportEmail, messageHandlers)
	if err != nil {
		log.Fatal(err)
	}
	defer handlerSet.Close()

	// Listen for incoming messages.
	err = handlerSet.Listen()
	if err != nil {
		log.Fatal(err)
	}

	// Spin until someone kills the process.
	spinner := make(chan int)
	<-spinner
}
