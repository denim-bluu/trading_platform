package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	pb "momentum-trading-platform/api/proto/data_service"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CacheEntry struct {
	Data       *pb.StockResponse
	Expiration time.Time
}

type Server struct {
	pb.UnimplementedDataServiceServer
	Logger     *log.Logger
	HttpClient HTTPClient
	Cache      sync.Map
	CacheTTL   time.Duration
	DB         *sql.DB
}

func NewServer(db *sql.DB) (*Server, error) {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	s := &Server{
		Logger:     logger,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
		CacheTTL:   15 * time.Minute,
		DB:         db,
	}

	if err := s.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return s, nil
}

func (s *Server) GetStockData(ctx context.Context, req *pb.StockRequest) (*pb.StockResponse, error) {
	s.Logger.WithFields(log.Fields{
		"symbol":     req.Symbol,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"interval":   req.Interval,
	}).Info("Received request for stock data")

	cacheKey := fmt.Sprintf("%s:%s:%s:%s", req.Symbol, req.StartDate, req.EndDate, req.Interval)

	// Check cache first
	if cachedData, found := s.getCachedData(cacheKey); found {
		s.Logger.Info("Returning cached data")
		return cachedData, nil
	}

	// Check database first
	data, err := s.getStockDataFromDB(req.Symbol, req.StartDate, req.EndDate)
	if err == nil {
		s.Logger.Info("Returning data from database")
		return data, nil
	}

	// If not in database, fetch from API
	data, err = s.fetchStockData(ctx, req.Symbol, req.StartDate, req.EndDate, req.Interval)
	if err != nil {
		return nil, err
	}

	// Store in database
	err = s.storeStockDataInDB(data)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to store data in database")
	}

	s.setCachedData(cacheKey, data)

	return data, nil
}

func (s *Server) GetBatchStockData(ctx context.Context, req *pb.BatchStockRequest) (*pb.BatchStockResponse, error) {
	s.Logger.WithFields(log.Fields{
		"symbols":    req.Symbols,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"interval":   req.Interval,
	}).Info("Received request for batch stock data")

	responses := make(map[string]*pb.StockResponse)
	errors := make(map[string]string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, symbol := range req.Symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()

			cacheKey := fmt.Sprintf("%s:%s:%s:%s", sym, req.StartDate, req.EndDate, req.Interval)
			if cachedData, found := s.getCachedData(cacheKey); found {
				mu.Lock()
				responses[sym] = cachedData
				mu.Unlock()
				s.Logger.WithField("symbol", sym).Info("Returning cached data")
				return
			}

			data, err := s.getStockDataFromDB(sym, req.StartDate, req.EndDate)
			if err == nil {
				mu.Lock()
				responses[sym] = data
				mu.Unlock()
				s.setCachedData(cacheKey, data)
				s.Logger.WithField("symbol", sym).Info("Returning data from database")
				return
			}

			time.Sleep(time.Millisecond * 100)
			data, err = s.fetchStockData(ctx, sym, req.StartDate, req.EndDate, req.Interval)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				s.Logger.WithError(err).WithField("symbol", sym).Error("Failed to fetch stock data")
				errors[sym] = err.Error()
			} else {
				responses[sym] = data
				go func() {
					if err := s.storeStockDataInDB(data); err != nil {
						s.Logger.WithError(err).Error("Failed to store data in database")
					}
					s.setCachedData(cacheKey, data)
				}()
			}
		}(symbol)
	}

	wg.Wait()

	return &pb.BatchStockResponse{
		StockData: responses,
		Errors:    errors,
	}, nil
}

func (s *Server) UpdateLatestData(ctx context.Context, req *pb.UpdateLatestDataRequest) (*pb.UpdateLatestDataResponse, error) {
	s.Logger.WithField("symbols", req.Symbols).Info("Updating latest data")

	err := s.updateDatabaseWithLatestData(req.Symbols)
	if err != nil {
		return &pb.UpdateLatestDataResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update latest data: %v", err),
		}, nil
	}

	return &pb.UpdateLatestDataResponse{
		Success: true,
		Message: "Successfully updated latest data",
	}, nil
}

func (s *Server) getCachedData(cacheKey string) (*pb.StockResponse, bool) {
	if entry, ok := s.Cache.Load(cacheKey); ok {
		cacheEntry := entry.(CacheEntry)
		if time.Now().Before(cacheEntry.Expiration) {
			return cacheEntry.Data, true
		}
		s.Cache.Delete(cacheKey)
	}
	return nil, false
}

func (s *Server) setCachedData(cacheKey string, data *pb.StockResponse) {
	s.Cache.Store(cacheKey, CacheEntry{
		Data:       data,
		Expiration: time.Now().Add(s.CacheTTL),
	})
}
