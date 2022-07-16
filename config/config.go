package config

import (
	"fmt"
	"os"
	"strconv"
)

type Telegram struct {
	Token string
}

type Config struct {
	Telegram             Telegram
	Secret               string
	StorageDSN           string
	NotificationDuration int
}

// New returns a new Config struct
func New() *Config {
	notDurationInt, _ := strconv.Atoi(getEnv("NOTIFICATION_DURATION", ""))

	return &Config{
		Telegram: Telegram{
			Token: getEnv("TELEGRAM_TOKEN", ""),
		},
		Secret:               getEnv("SECRET", ""),
		NotificationDuration: notDurationInt,
		StorageDSN:           fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", getEnv("POSTGRES_USER", ""), getEnv("POSTGRES_PASSWORD", ""), getEnv("POSTGRES_DBNAME", ""), getEnv("POSTGRES_SSLMODE", "")),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
