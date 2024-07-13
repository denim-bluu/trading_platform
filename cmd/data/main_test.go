// cmd/data/main_test.go
package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	pb "momentum-trading-platform/api/proto/data_service"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/time/rate"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetStockData(t *testing.T) {
	mockClient := new(MockHTTPClient)

	// Create a mock response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`{
            "chart": {
                "result": [{
                    "meta": {"symbol": "AAPL"},
                    "timestamp": [1625097600],
                    "indicators": {
                        "quote": [{
                            "open": [100],
                            "high": [101],
                            "low": [99],
                            "close": [100.5],
                            "volume": [1000000]
                        }],
                        "adjclose": [{
                            "adjclose": [100.5]
                        }]
                    }
                }]
            }
        }`)),
	}

	// Set up expectations
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	s := &server{
		logger:      logrus.New(),
		rateLimiter: rate.NewLimiter(rate.Every(time.Second), 2),
		httpClient:  mockClient,
	}

	ctx := context.Background()
	req := &pb.StockRequest{
		Symbol:    "AAPL",
		StartDate: "1625097600",
		EndDate:   "1625184000",
		Interval:  "1d",
	}

	resp, err := s.GetStockData(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "AAPL", resp.Symbol)
	assert.Len(t, resp.DataPoints, 1)
	assert.Equal(t, int64(1625097600), resp.DataPoints[0].Timestamp)
	assert.Equal(t, 100.0, resp.DataPoints[0].Open)
	assert.Equal(t, 101.0, resp.DataPoints[0].High)
	assert.Equal(t, 99.0, resp.DataPoints[0].Low)
	assert.Equal(t, 100.5, resp.DataPoints[0].Close)
	assert.Equal(t, 100.5, resp.DataPoints[0].AdjustedClose)
	assert.Equal(t, int64(1000000), resp.DataPoints[0].Volume)

	// Verify that the expectations were met
	mockClient.AssertExpectations(t)
}

func TestGetBatchStockData(t *testing.T) {
	mockClient := new(MockHTTPClient)

	// Create mock responses for AAPL and MSFT from Yahoo Finance
	mockResponseAAPL := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`{
            "chart": {
                "result": [{
                    "meta": {"symbol": "AAPL"},
                    "timestamp": [1625097600],
                    "indicators": {
                        "quote": [{
                            "open": [100],
                            "high": [101],
                            "low": [99],
                            "close": [100.5],
                            "volume": [1000000]
                        }],
                        "adjclose": [{
                            "adjclose": [100.5]
                        }]
                    }
                }]
            }
        }`)),
	}

	mockResponseMSFT := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`{
            "chart": {
                "result": [{
                    "meta": {"symbol": "MSFT"},
                    "timestamp": [1625097600],
                    "indicators": {
                        "quote": [{
                            "open": [200],
                            "high": [201],
                            "low": [199],
                            "close": [200.5],
                            "volume": [2000000]
                        }],
                        "adjclose": [{
                            "adjclose": [200.5]
                        }]
                    }
                }]
            }
        }`)),
	}

	// Set up expectations with improved matcher
	// This matcher will check if the URL is correct
	// for the given stock symbol
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		url := req.URL.String()
		return url == "https://query1.finance.yahoo.com/v8/finance/chart/AAPL?period1=1625097600&period2=1625184000&interval=1d"
	})).Return(mockResponseAAPL, nil)

	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		url := req.URL.String()
		return url == "https://query1.finance.yahoo.com/v8/finance/chart/MSFT?period1=1625097600&period2=1625184000&interval=1d"
	})).Return(mockResponseMSFT, nil)

	s := &server{
		logger:      logrus.New(),
		rateLimiter: rate.NewLimiter(rate.Every(time.Second), 2),
		httpClient:  mockClient,
	}

	ctx := context.Background()
	req := &pb.BatchStockRequest{
		Symbols:   []string{"AAPL", "MSFT"},
		StartDate: "1625097600",
		EndDate:   "1625184000",
		Interval:  "1d",
	}

	resp, err := s.GetBatchStockData(ctx, req)

	assert.NoError(t, err)
	assert.Contains(t, resp.StockData, "AAPL")
	assert.Contains(t, resp.StockData, "MSFT")

	assert.Equal(t, "AAPL", resp.StockData["AAPL"].Symbol)
	assert.Len(t, resp.StockData["AAPL"].DataPoints, 1)
	assert.Equal(t, int64(1625097600), resp.StockData["AAPL"].DataPoints[0].Timestamp)
	assert.Equal(t, 100.0, resp.StockData["AAPL"].DataPoints[0].Open)
	assert.Equal(t, 101.0, resp.StockData["AAPL"].DataPoints[0].High)
	assert.Equal(t, 99.0, resp.StockData["AAPL"].DataPoints[0].Low)
	assert.Equal(t, 100.5, resp.StockData["AAPL"].DataPoints[0].Close)
	assert.Equal(t, 100.5, resp.StockData["AAPL"].DataPoints[0].AdjustedClose)
	assert.Equal(t, int64(1000000), resp.StockData["AAPL"].DataPoints[0].Volume)

	assert.Equal(t, "MSFT", resp.StockData["MSFT"].Symbol)
	assert.Len(t, resp.StockData["MSFT"].DataPoints, 1)
	assert.Equal(t, int64(1625097600), resp.StockData["MSFT"].DataPoints[0].Timestamp)
	assert.Equal(t, 200.0, resp.StockData["MSFT"].DataPoints[0].Open)
	assert.Equal(t, 201.0, resp.StockData["MSFT"].DataPoints[0].High)
	assert.Equal(t, 199.0, resp.StockData["MSFT"].DataPoints[0].Low)
	assert.Equal(t, 200.5, resp.StockData["MSFT"].DataPoints[0].Close)
	assert.Equal(t, 200.5, resp.StockData["MSFT"].DataPoints[0].AdjustedClose)
	assert.Equal(t, int64(2000000), resp.StockData["MSFT"].DataPoints[0].Volume)

	// Verify that the expectations were met
	mockClient.AssertExpectations(t)
}
