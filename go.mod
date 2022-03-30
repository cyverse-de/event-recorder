module github.com/cyverse-de/event-recorder

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/DavidGamba/go-getoptions v0.20.0
	github.com/Masterminds/squirrel v1.4.0
	github.com/cyverse-de/configurate v0.0.0-20200527185205-4e1e92866cee
	github.com/cyverse-de/dbutil v1.0.1
	github.com/cyverse-de/messaging/v9 v9.1.1
	github.com/lib/pq v1.10.4
	github.com/mcnijman/go-emailaddress v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.7.0 // indirect
	github.com/streadway/amqp v1.0.1-0.20200716223359-e6b33f460591
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.6.1
	go.opentelemetry.io/otel/exporters/jaeger v1.6.1
	go.opentelemetry.io/otel/sdk v1.6.1
)
