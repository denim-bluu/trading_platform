package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	pb "trading_platform/strategy/proto"
)

type DataFetcher interface {
	FetchData(symbol, startDate, endDate string) (*pb.HistoricalData, error)
}

type httpDataFetcher struct {
	aggregatorURL string
}

func NewHTTPDataFetcher(aggregatorURL string) DataFetcher {
	return &httpDataFetcher{aggregatorURL: aggregatorURL}
}

func (f *httpDataFetcher) FetchData(symbol, startDate, endDate string) (*pb.HistoricalData, error) {
	// Calculate the buffer start date (200 days before start date)
	bufferStartDate := calculateBufferStartDate(startDate, 200)

	url := fmt.Sprintf("%s/get_historical_data", f.aggregatorURL)
	reqBody := map[string]string{
		"symbol":     symbol,
		"start_date": bufferStartDate,
		"end_date":   endDate,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Log the raw response body
	log.Printf("Raw response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}

	var data pb.HistoricalData
	if err := json.Unmarshal(body, &data); err != nil {
		// Log the unmarshalling error
		log.Printf("Failed to unmarshal response: %v", err)
		return nil, err
	}
	return &data, nil
}

// calculateBufferStartDate calculates a start date 200 days before the provided start date
func calculateBufferStartDate(startDate string, bufferDays int) string {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, startDate)
	bufferStartDate := t.AddDate(0, 0, -bufferDays)
	return bufferStartDate.Format(layout)
}
