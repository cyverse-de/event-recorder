package main

import (
	"fmt"
	"os"

	"github.com/DavidGamba/go-getoptions"
	"github.com/cyverse-de/configurate"
	"github.com/cyverse-de/logcabin"
)

// commandLineOptionValues represents the values of the command-line options that were passed on the command line when
// this service was invoked.
type commandLineOptionValues struct {
	Config string
}

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
	amqpURI := cfg.GetString("amqp.uri")
	amqpExchangeName := cfg.GetString("amqp.exchange.name")
	amqpExchangeType := cfg.GetString("amqp.exchange.type")

	fmt.Printf("%s\n", amqpURI)
	fmt.Printf("%s\n", amqpExchangeName)
	fmt.Printf("%s\n", amqpExchangeType)
}
