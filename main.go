package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/DavidGamba/go-getoptions"
	"github.com/cyverse-de/configurate"
	"github.com/cyverse-de/event-recorder/common"
	"github.com/cyverse-de/event-recorder/db"
	"github.com/cyverse-de/event-recorder/handlers"
	"github.com/cyverse-de/event-recorder/handlerset"
	"github.com/cyverse-de/logcabin"
)

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
	}

	return optionValues
}

func main() {
	// Parse the command-line.
	optionValues := parseCommandLine()

	// Initialize logging.
	logcabin.Init("event-recorder", "event-recorder")

	// Read in the configuration file.
	cfg, err := configurate.InitDefaults(optionValues.Config, configurate.JobServicesDefaults)
	if err != nil {
		logcabin.Error.Fatal(err)
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
		logcabin.Error.Fatal(err)
	}
	defer db.Close()

	// Get the email address to use for support requests.
	supportEmail := cfg.GetString("email.request")

	// Initialize the message handlers.
	messageHandlers, err := handlers.InitMessageHandlers(db, amqpSettings)
	if err != nil {
		logcabin.Error.Fatal(err)
	}

	// Create the message handler set.
	handlerSet, err := handlerset.New(amqpSettings, supportEmail, messageHandlers)
	if err != nil {
		logcabin.Error.Fatal(err)
	}
	defer handlerSet.Close()

	// Listen for incoming messages.
	err = handlerSet.Listen()
	if err != nil {
		logcabin.Error.Fatal(err)
	}

	// Spin until someone kills the process.
	spinner := make(chan int)
	<-spinner
}
