package twitch

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Username        string
	OAuthToken      string
	ClientID        string
	Channel         string
	DbPath          string
	ReconnectConfig ReconnectConfig
}

type ReconnectConfig struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		Username:   getEnvOrDefault("TWITCH_USERNAME", "kitsunebot"),
		OAuthToken: getEnvOrDefault("TWITCH_OAUTH", "oauth:placeholder"),
		ClientID:   getEnvOrDefault("TWITCH_CLIENT_ID", "client_id_placeholder"),
		Channel:    getEnvOrDefault("TWITCH_CHANNEL", "kitsunebot"),
		DbPath:     getEnvOrDefault("DB_PATH", "./kitsunebot.db"),
		ReconnectConfig: ReconnectConfig{
			InitialDelay: getEnvDurationOrDefault("RECONNECT_INITIAL_DELAY", 1*time.Second),
			MaxDelay:     getEnvDurationOrDefault("RECONNECT_MAX_DELAY", 30*time.Second),
			Multiplier:   getEnvFloatOrDefault("RECONNECT_MULTIPLIER", 2.0),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		var f float64
		if _, err := fmt.Sscanf(value, "%f", &f); err == nil {
			return f
		}
	}
	return defaultValue
}
