// utils/calculations.go
package utils

import (
	"fmt"
	"strconv"
	"time"
)

func ParseAndFormatTimestamp(timestampStr string) (string, error) {
	// Parse int64 from string
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Format the timestamp into a date string
	date := time.Unix(timestamp, 0).Format("2006-01-02:15:04:05")
	return date, nil
}
