package config

import "os"

type Config struct {
	AGGREGATOR_URL string
	Port           string
}

func Load() (*Config, error) {
	return &Config{
		AGGREGATOR_URL: getEnv("AGGREGATOR_URL", "http://localhost:8080"),
		Port:           getEnv("PORT", "8081"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
