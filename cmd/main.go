package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sword-challenge/internal/api"
	"github.com/sword-challenge/internal/environment"
	"github.com/sword-challenge/internal/service/consumeNotify"
	"github.com/sword-challenge/pkg/rabbitmq"
	"github.com/sword-challenge/pkg/rabbitmq/models"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	godotenv.Load()
	environment.CheckEnvVars()

	// configure RabbitMQ
	publisher, consumer := setupRabbit()

	// Start consume service
	go consumeEvents(consumer)

	// configure DB
	db := setupDatabase()

	// Create an instance of echo and register all the routes
	e := setupServer()
	a := api.Start(e, db, publisher, consumer)
	a.RegisterRoutes()

	log.Info().Msgf("HTTP Server on port %v", os.Getenv(environment.Port))
	e.Logger.Fatal(e.Start(":" + os.Getenv(environment.Port)))
}

func setupServer() *echo.Echo {
	// create an instance of Echo and update the default config
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return e
}

func setupRabbit() (rabbitmq.Publisher, rabbitmq.Consumer) {
	client, err := rabbitmq.New(models.Credential{
		Host:     os.Getenv(environment.RabbitHost),
		User:     os.Getenv(environment.RabbitUsername),
		Password: os.Getenv(environment.RabbitPassword),
	})
	if err != nil {
		log.Error().Err(err).Msg("Error trying to create a new RabbitMQ client")
	}

	publisher, err := client.NewPublisher(&models.QueueArgs{
		Name: os.Getenv(environment.TaskQueue),
	})
	if err != nil {
		log.Error().Err(err).Msg("Error trying to create a new RabbitMQ Publisher client")
	}

	consumer, err := client.NewConsumer(os.Getenv(environment.TaskQueue))
	if err != nil {
		log.Error().Err(err).Msg("Error trying to create a new RabbitMQ Consumer client")
	}

	return publisher, consumer
}

func setupDatabase() *sqlx.DB {
	databaseUrl := os.Getenv(environment.DBUser) + ":" + os.Getenv(environment.DBPassword) + "@tcp(" + os.Getenv(environment.DBHost) + ")/" + os.Getenv(environment.DBName) + "?multiStatements=true&parseTime=true"

	db, err := sqlx.Open("mysql", databaseUrl)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to the database. error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.Fatal().Msgf("Failed to create a new instance of mysql. error: %v", err)
	}

	// m, err := migrate.New("file://db/migrations", databaseUrl)
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		os.Getenv("DB_NAME"), driver)
	if err != nil {
		log.Fatal().Msgf("Failed to create a new instance of go-migrate. error: %v", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("Failed to migrate database to latest version. error %v", err)
	}

	return db
}

func consumeEvents(consumer rabbitmq.Consumer) {
	// Start to listening queue
	log.Info().Msg("Started to listening events")
	service := consumeNotify.New()

	ev := models.ConsumerEvent{
		Handler:            service.ConsumeEvent,
		RetryMessagePeriod: 0,
		QueueName:          environment.TaskQueue,
	}

	if err := consumer.SubscribeEvents(context.Background(), ev); err != nil {
		log.Error().Err(err).Msg("Error on Subscribe events")
		panic(fmt.Sprintln("Error on subscribe events", err))
	}
}
