package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

// AppConfig specifies high level configuration settings for the app
type AppConfig struct {
	Port       int
	DBFile     string
	SessionKey string
}

// ParseAppConfig parses flags and/or env vars returning an AppConfig instance
func ParseAppConfig() AppConfig {
	log.Print("\n\n  _______          _           _ \n" +
		" |__   __|        (_)         | |\n" +
		"    | | ___  _ __  _  ___ __ _| |\n" +
		"    | |/ _ \\| '_ \\| |/ __/ _` | |\n" +
		"    | | (_) | |_) | | (_| (_| | |\n" +
		"    |_|\\___/| .__/|_|\\___\\__,_|_|\n" +
		"            | |                  \n" +
		"            |_|\n\n")

	port := flag.Int("p", envOrInt("PORT", 8000), "port for app on run on")
	dbFile := flag.String("db-file", envOrString("DB_FILE", "topical.db"), "path to sqlite db file")
	sessionKey := flag.String("session-key", envOrString("SESSION_KEY", "not-set"), "session key for cookie store")

	flag.Parse()

	return AppConfig{*port, *dbFile, *sessionKey}
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
