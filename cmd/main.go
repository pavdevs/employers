package main

import (
	"employer.dev/internal/api"
	"employer.dev/internal/database"
	employerrepository "employer.dev/internal/repository/employer"
	"employer.dev/internal/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = logrus.New()
	loadEnv()
}

// @title Employers service
// @version 0.0.1
// @description Web-server for employers

// @host localhost:8000
// @BasePath /
func main() {
	dbService := prepareDatabase()

	if err := dbService.Connect(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	rep := employerrepository.NewEmployerRepository(dbService, logger)
	empAPI := api.NewEmployerAPI(rep)
	webServer := server.NewServer(
		server.NewConfig(os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
		empAPI,
		logger,
	)

	if err := webServer.Start(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func prepareDatabase() *database.Database {
	return database.NewDatabase(database.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DBName:   os.Getenv("DATABASE_NAME"),
		SSLMode:  os.Getenv("DATABASE_SSLMODE"),
	}, logger)
}

func loadEnv() {
	if err := godotenv.Load("cmd/config.env"); err != nil {
		logger.Error("No .env file found")
	}
}
