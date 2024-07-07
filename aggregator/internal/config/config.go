package config

import "os"

type Config struct {
	YahooFinanceAPIURL string
	Port               string
}

func Load() (*Config, error) {
	return &Config{
		YahooFinanceAPIURL: getEnv("YAHOO_FINANCE_API_URL", "https://query1.finance.yahoo.com"),
		Port:               getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
