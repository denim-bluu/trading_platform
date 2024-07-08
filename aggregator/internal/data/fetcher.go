package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pb "trading_platform/aggregator/proto"
)

type DataFetcher interface {
	FetchHistoricalData(symbol string, startDate string, endDate string) ([]*pb.DataPoint, error)
	FetchLiveData(symbol string) (*pb.DataPoint, error)
}

type YahooFinanceFetcher struct {
	apiURL string
}

func NewYahooFinanceFetcher(apiURL string) DataFetcher {
	return &YahooFinanceFetcher{apiURL: apiURL}
}

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open     []float64 `json:"open"`
					High     []float64 `json:"high"`
					Low      []float64 `json:"low"`
					Close    []float64 `json:"close"`
					Adjclose []float64 `json:"adjclose"`
					Volume   []float64 `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func (yf *YahooFinanceFetcher) FetchHistoricalData(symbol string, startDate string, endDate string) ([]*pb.DataPoint, error) {
	url := fmt.Sprintf("%s/v8/finance/chart/%s?period1=%d&period2=%d&interval=1d", yf.apiURL, symbol, toUnixTime(startDate), toUnixTime(endDate))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data yahooResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	quotes := data.Chart.Result[0].Indicators.Quote[0]
	timestamps := data.Chart.Result[0].Timestamp

	result := make([]*pb.DataPoint, len(timestamps))
	for i := range timestamps {
		result[i] = &pb.DataPoint{
			Symbol:    symbol,
			Timestamp: time.Unix(timestamps[i], 0).Format("2006-01-02"),
			Open:      quotes.Open[i],
			High:      quotes.High[i],
			Low:       quotes.Low[i],
			Close:     quotes.Close[i],
			Adjclose:  quotes.Adjclose[i],
			Volume:    quotes.Volume[i],
		}
	}
	return result, nil
}

func (yf *YahooFinanceFetcher) FetchLiveData(symbol string) (*pb.DataPoint, error) {
	url := fmt.Sprintf("%s/v8/finance/chart/%s?interval=1m", yf.apiURL, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data yahooResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	quotes := data.Chart.Result[0].Indicators.Quote[0]
	timestamps := data.Chart.Result[0].Timestamp
	latestIndex := len(timestamps) - 1

	liveData := &pb.DataPoint{
		Timestamp: time.Unix(timestamps[latestIndex], 0).Format("2006-01-02 15:04:05"),
		Open:      quotes.Open[latestIndex],
		High:      quotes.High[latestIndex],
		Low:       quotes.Low[latestIndex],
		Close:     quotes.Close[latestIndex],
		Adjclose:  quotes.Adjclose[latestIndex],
		Volume:    quotes.Volume[latestIndex],
	}

	return liveData, nil
}

func toUnixTime(dateStr string) int64 {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dateStr)
	return t.Unix()
}
