// cmd/data/integration_test.go
package main

import (
	"context"
	"testing"
	"time"

	pb "momentum-trading-platform/api/proto/data_service"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestIntegrationGetStockData(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	startDate := "1640995200" // 2022-01-01
	endDate := "1672531199"   // 2022-12-31

	req := &pb.StockRequest{
		Symbol:    "AAPL",
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "1mo",
	}

	resp, err := client.GetStockData(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "AAPL", resp.Symbol)
	assert.NotEmpty(t, resp.DataPoints)

	// Check if we have approximately 12 data points (monthly data for a year)
	assert.True(t, len(resp.DataPoints) >= 11 && len(resp.DataPoints) <= 13,
		"Expected around 12 data points, got %d", len(resp.DataPoints))

	// Check if the data points are within the specified date range
	firstPoint := resp.DataPoints[0]
	lastPoint := resp.DataPoints[len(resp.DataPoints)-1]
	assert.True(t, firstPoint.Timestamp >= int64(1640995200),
		"First data point is before start date")
	assert.True(t, lastPoint.Timestamp <= int64(1672531199),
		"Last data point is after end date")

	// Basic sanity checks on the data
	for _, point := range resp.DataPoints {
		assert.True(t, point.Open > 0, "Open price should be positive")
		assert.True(t, point.High >= point.Low, "High should be >= Low")
		assert.True(t, point.Close > 0, "Close price should be positive")
		assert.True(t, point.AdjustedClose > 0, "Adjusted close should be positive")
		assert.True(t, point.Volume >= 0, "Volume should be non-negative")
	}
}

func TestIntegrationGetBatchStockData(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startDate := "1640995200" // 2022-01-01
	endDate := "1672531199"   // 2022-12-31

	req := &pb.BatchStockRequest{
		Symbols:   []string{"AAPL", "GOOGL", "MSFT"},
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  "1mo",
	}

	resp, err := client.GetBatchStockData(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, resp.StockData, 3)
	assert.Contains(t, resp.StockData, "AAPL")
	assert.Contains(t, resp.StockData, "GOOGL")
	assert.Contains(t, resp.StockData, "MSFT")

	for symbol, stockResp := range resp.StockData {
		assert.Equal(t, symbol, stockResp.Symbol)
		assert.NotEmpty(t, stockResp.DataPoints)

		// Check if we have approximately 12 data points (monthly data for a year)
		assert.True(t, len(stockResp.DataPoints) >= 11 && len(stockResp.DataPoints) <= 13,
			"Expected around 12 data points for %s, got %d", symbol, len(stockResp.DataPoints))

		// Check if the data points are within the specified date range
		firstPoint := stockResp.DataPoints[0]
		lastPoint := stockResp.DataPoints[len(stockResp.DataPoints)-1]
		assert.True(t, firstPoint.Timestamp >= int64(1640995200),
			"First data point for %s is before start date", symbol)
		assert.True(t, lastPoint.Timestamp <= int64(1672531199),
			"Last data point for %s is after end date", symbol)

		// Basic sanity checks on the data
		for _, point := range stockResp.DataPoints {
			assert.True(t, point.Open > 0, "Open price should be positive for %s", symbol)
			assert.True(t, point.High >= point.Low, "High should be >= Low for %s", symbol)
			assert.True(t, point.Close > 0, "Close price should be positive for %s", symbol)
			assert.True(t, point.AdjustedClose > 0, "Adjusted close should be positive for %s", symbol)
			assert.True(t, point.Volume >= 0, "Volume should be non-negative for %s", symbol)
		}
	}
}
