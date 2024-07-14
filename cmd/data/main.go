// cmd/data/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	pb "momentum-trading-platform/api/proto/data_service"
	"momentum-trading-platform/utils"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type server struct {
	pb.UnimplementedDataServiceServer
	logger      *log.Logger
	rateLimiter *rate.Limiter
	httpClient  HTTPClient
}

func newServer() *server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return &server{
		logger:      logger,
		rateLimiter: rate.NewLimiter(rate.Every(time.Second), 2),
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *server) GetStockData(ctx context.Context, req *pb.StockRequest) (*pb.StockResponse, error) {
	startTimestamp, _ := utils.ParseAndFormatTimestamp(req.StartDate)
	endTimestamp, _ := utils.ParseAndFormatTimestamp(req.EndDate)
	s.logger.WithFields(log.Fields{
		"symbol":   req.Symbol,
		"start":    startTimestamp,
		"end":      endTimestamp,
		"interval": req.Interval,
	}).Info("Received request for stock data")

	if err := s.rateLimiter.Wait(ctx); err != nil {
		s.logger.WithError(err).Error("Rate limit exceeded")
		return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}

	return s.fetchStockData(ctx, req.Symbol, req.StartDate, req.EndDate, req.Interval)
}

func (s *server) GetBatchStockData(ctx context.Context, req *pb.BatchStockRequest) (*pb.BatchStockResponse, error) {
	s.logger.WithFields(log.Fields{
		"symbols":  req.Symbols,
		"start":    req.StartDate,
		"end":      req.EndDate,
		"interval": req.Interval,
	}).Info("Received request for batch stock data")

	responses := make(map[string]*pb.StockResponse)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, symbol := range req.Symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			// Apply rate limiting
			if err := s.rateLimiter.Wait(ctx); err != nil {
				s.logger.WithError(err).WithField("symbol", sym).Error("Rate limit exceeded")
				return
			}

			resp, err := s.fetchStockData(ctx, sym, req.StartDate, req.EndDate, req.Interval)
			if err != nil {
				s.logger.WithError(err).WithField("symbol", sym).Error("Failed to fetch stock data")
				return
			}

			mu.Lock()
			responses[sym] = resp
			mu.Unlock()
		}(symbol)
	}

	wg.Wait()

	return &pb.BatchStockResponse{StockData: responses}, nil
}

func (s *server) fetchStockData(ctx context.Context, symbol, startDate, endDate, interval string) (*pb.StockResponse, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?period1=%s&period2=%s&interval=%s",
		symbol, startDate, endDate, interval)

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create request")
		return nil, status.Errorf(codes.Internal, "failed to create request: %v", err)
	}

	// Get response from Yahoo Finance
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		s.logger.WithError(err).Error("Failed to fetch data from Yahoo Finance")
		return nil, status.Errorf(codes.Internal, "failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.WithField("status", resp.StatusCode).Error("Received non-200 response from Yahoo Finance")
		return nil, status.Errorf(codes.Internal, "received non-200 response: %d", resp.StatusCode)
	}

	var yahooResp yahooFinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooResp); err != nil {
		s.logger.WithError(err).Error("Failed to decode response from Yahoo Finance")
		return nil, status.Errorf(codes.Internal, "failed to decode response: %v", err)
	}

	if len(yahooResp.Chart.Result) == 0 {
		s.logger.WithField("symbol", symbol).Warn("No data found for symbol")
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

	return &pb.StockResponse{
		Symbol:     symbol,
		DataPoints: dataPoints,
	}, nil
}

func main() {
	s := newServer()

	lis, err := net.Listen("tcp", utils.GetEnv("GRPC_ADDR", ":50051"))
	if err != nil {
		s.logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, s)

	s.logger.WithField("address", lis.Addr().String()).Info("Data service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.WithError(err).Fatal("Failed to serve")
	}
}
