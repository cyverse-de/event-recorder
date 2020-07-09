module github.com/cyverse-de/event-recorder

go 1.14

require (
	github.com/DavidGamba/go-getoptions v0.20.0
	github.com/Masterminds/squirrel v1.4.0
	github.com/cyverse-de/configurate v0.0.0-20200527185205-4e1e92866cee
	github.com/cyverse-de/dbutil v0.0.0-20200527185309-2b32eb41f45e
	github.com/cyverse-de/logcabin v0.0.0-20200527185931-2ccd850e39ad
	github.com/cyverse-de/messaging v6.0.0+incompatible
	github.com/cyverse-de/model v0.0.0-20200527190032-dc1e3a7c2ccd // indirect
	github.com/fatih/structs v1.1.0
	github.com/lib/pq v1.7.0
	github.com/mcnijman/go-emailaddress v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.7.0 // indirect
	github.com/streadway/amqp v1.0.0
	gopkg.in/cyverse-de/model.v4 v4.0.0-20191010001558-736b5a572acd // indirect
)

replace github.com/cyverse-de/messaging => /Users/sarahr/src/de/go/src/github.com/cyverse-de/messaging
