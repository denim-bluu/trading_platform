package config

import (
	"os"
	"strconv"
)

type Config struct {
	YahooFinanceAPIURL string
	Port               string
	DataDir            string
	CacheTTL           int
}

func Load() (*Config, error) {
	return &Config{
		YahooFinanceAPIURL: getEnv("YAHOO_FINANCE_API_URL", "https://query1.finance.yahoo.com"),
		Port:               getEnv("PORT", "8080"),
		DataDir:            getEnv("DATA_DIR", "app/data"),
		CacheTTL:           parseInt(getEnv("CACHE_TTL", "24")),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}
