package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	AppPort    int
	DBDriver   string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	SessionKey string
}

func ParseAppConfig() AppConfig {
	log.Print("\n\n  _______          _           _ \n" +
		" |__   __|        (_)         | |\n" +
		"    | | ___  _ __  _  ___ __ _| |\n" +
		"    | |/ _ \\| '_ \\| |/ __/ _` | |\n" +
		"    | | (_) | |_) | | (_| (_| | |\n" +
		"    |_|\\___/| .__/|_|\\___\\__,_|_|\n" +
		"            | |                  \n" +
		"            |_|\n\n")

	port := flag.Int("p", envOrInt("APP_PORT", 8080), "port for app on run on")
	dbDriver := flag.String("db-driver", envOrString("DB_DRIVER", "postgres"), "db driver to use for app")
	dbHost := flag.String("db-host", envOrString("DB_HOST", "localhost"), "host for Postgres DB")
	dbPort := flag.Int("db-port", envOrInt("DB_PORT", 5432), "port for Postgres DB")
	dbUser := flag.String("db-user", envOrString("DB_USER", "topical"), "port for Postgres DB")
	dbPassword := flag.String("db-password", envOrString("DB_PASSWORD", "not-set"), "password for Postgres DB")
	dbName := flag.String("db-name", envOrString("DB_NAME", "topical"), "name for Postgres DB")
	sslMode := flag.String("db-ssl-mode", envOrString("DB_SSL_MODE", "disable"), "whether to enable or disable connecting to DB in SSL mode")
	sessionKey := flag.String("session-key", envOrString("SESSION_KEY", "not-set"), "session key for cookie store")

	flag.Parse()

	return AppConfig{*port, *dbDriver, *dbHost, *dbPort, *dbUser, *dbPassword, *dbName, *sslMode, *sessionKey}
}

func envOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func envOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("envOrInt[%s]: %v", key, err)
		}
		return value
	}
	return defaultVal
}
