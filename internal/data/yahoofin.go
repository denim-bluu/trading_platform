package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "momentum-trading-platform/api/proto/data_service"
	"momentum-trading-platform/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type yahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol string `json:"symbol"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
				Adjclose []struct {
					Adjclose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func (s *Server) fetchStockData(ctx context.Context, symbol, startDate, endDate, interval string) (*pb.StockResponse, error) {
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", symbol, startDate, endDate, interval)

	// Check cache first
	if cachedData, found := s.getCachedData(cacheKey); found {
		s.Logger.WithField("symbol", symbol).Info("Returning cached data")
		return cachedData, nil
	}

	// If not in cache, fetch from Yahoo Finance
	startDateUnix, err := utils.ConvertDateStrToUnixTimestamp(startDate, "2006-01-02")
	if err != nil {
		s.Logger.WithError(err).Error("Failed to convert start date")
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert start date: %v", err)
	}

	endDateUnix, err := utils.ConvertDateStrToUnixTimestamp(endDate, "2006-01-02")
	if err != nil {
		s.Logger.WithError(err).Error("Failed to convert end date")
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert end date: %v", err)
	}

	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?period1=%s&period2=%s&interval=%s",
		symbol, startDateUnix, endDateUnix, interval)

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create request")
		return nil, status.Errorf(codes.Internal, "failed to create request: %v", err)
	}

	// Get response from Yahoo Finance
	resp, err := s.HttpClient.Do(httpReq)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch data from Yahoo Finance")
		return nil, status.Errorf(codes.Internal, "failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.Logger.WithField("status", resp.StatusCode).Error("Received non-200 response from Yahoo Finance")
		return nil, status.Errorf(codes.Internal, "received non-200 response: %d", resp.StatusCode)
	}

	var yahooResp yahooFinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooResp); err != nil {
		s.Logger.WithError(err).Error("Failed to decode response from Yahoo Finance")
		return nil, status.Errorf(codes.Internal, "failed to decode response: %v", err)
	}

	if len(yahooResp.Chart.Result) == 0 {
		s.Logger.WithField("symbol", symbol).Warn("No data found for symbol")
		return nil, status.Errorf(codes.NotFound, "no data found for symbol: %s", symbol)
	}

	result := yahooResp.Chart.Result[0]
	dataPoints := make([]*pb.StockDataPoint, len(result.Timestamp))
	for i, ts := range result.Timestamp {
		dataPoints[i] = &pb.StockDataPoint{
			Timestamp:     ts,
			Open:          result.Indicators.Quote[0].Open[i],
			High:          result.Indicators.Quote[0].High[i],
			Low:           result.Indicators.Quote[0].Low[i],
			Close:         result.Indicators.Quote[0].Close[i],
			AdjustedClose: result.Indicators.Adjclose[0].Adjclose[i],
			Volume:        result.Indicators.Quote[0].Volume[i],
		}
	}

	stockResponse := &pb.StockResponse{
		Symbol:     symbol,
		DataPoints: dataPoints,
	}

	// Cache the fetched data
	s.setCachedData(cacheKey, stockResponse)

	return stockResponse, nil
}
