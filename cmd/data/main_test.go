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
