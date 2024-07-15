// utils/calculations.go
package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/piquette/finance-go/datetime"
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

func ParseDateStrToDateTime(dateStr string) (*datetime.Datetime, error) {
	const layout = "2006-01-02" // Change this layout to match your date format
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}
	return datetime.New(&t), nil
}

func ConvertDateStrToUnixTimestamp(dateStr, layout string) (string, error) {
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", t.Unix()), nil
}
