package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sidartaoss/fullcycle/encoder/application/services"
	"github.com/sidartaoss/fullcycle/encoder/framework/database"
	"github.com/sidartaoss/fullcycle/encoder/framework/queue"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file, %v", err)
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("error parsing boolean var; %v", err)
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("error parsing boolean var; %v", err)
	}

	db.AutoMigrateDb = autoMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("error connecting to DB, %v", err)
	}

	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()

	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(
		dbConnection,
		rabbitMQ,
		jobReturnChannel,
		messageChannel,
	)

	jobManager.Start(ch)
}
