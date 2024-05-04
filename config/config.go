package config

import (
	"enigmanations/cats-social/pkg/env"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	AppPort    int
	AppHost    string
	DBHost     string
	DBUsername string
	DBPass     string
	DBName     string
	DBPort     int
	DBParams   string
	Salt       int
}

func GetConfig() *Configuration {
	// Load env
	initEnv()

	appPort, err := env.GetEnvInt("APP_PORT")
	if err != nil {
		return nil
	}

	appHost, err := env.GetEnv("APP_HOST")
	if err != nil {
		return nil
	}

	dbHost, err := env.GetEnv("DB_HOST")
	if err != nil {
		return nil
	}

	dbUsername, err := env.GetEnv("DB_USERNAME")
	if err != nil {
		return nil
	}

	dbPass, err := env.GetEnv("DB_PASSWORD")
	if err != nil {
		return nil
	}

	dbName, err := env.GetEnv("DB_NAME")
	if err != nil {
		return nil
	}

	dbPort, err := env.GetEnvInt("DB_PORT")
	if err != nil {
		return nil
	}

	dbParams, err := env.GetEnv("DB_PARAMS")
	if err != nil {
		return nil
	}

	salt, err := env.GetEnvInt("BCRYPT_SALT")
	if err != nil {
		return nil
	}

	return &Configuration{
		AppPort:    appPort,
		AppHost:    appHost,
		DBHost:     dbHost,
		DBUsername: dbUsername,
		DBPass:     dbPass,
		DBName:     dbName,
		DBPort:     dbPort,
		DBParams:   dbParams,
		Salt:       salt,
	}
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env %v\n", err)
		os.Exit(1)
	}
}
