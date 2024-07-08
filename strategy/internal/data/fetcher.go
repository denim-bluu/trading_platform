package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pb "trading_platform/strategy/proto"
)

type DataFetcher interface {
	FetchData(symbol, startDate, endDate string) (*pb.HistoricalData, error)
	FetchIndexValue(index string) (float64, error)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}

	var data pb.HistoricalData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *httpDataFetcher) FetchIndexValue(index string) (float64, error) {
	url := fmt.Sprintf("%s/get_historical_data", f.aggregatorURL)
	reqBody := map[string]string{
		"symbol": index,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch index value: status code %d", resp.StatusCode)
	}

	var response pb.HistoricalData
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	if len(response.DataPoints) == 0 {
		return 0, fmt.Errorf("no data points found for index: %s", index)
	}
	latestDataPoint := response.DataPoints[len(response.DataPoints)-1]
	return latestDataPoint.Close, nil
}

// calculateBufferStartDate calculates a start date 200 days before the provided start date
func calculateBufferStartDate(startDate string, bufferDays int) string {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, startDate)
	bufferStartDate := t.AddDate(0, 0, -bufferDays)
	return bufferStartDate.Format(layout)
}
