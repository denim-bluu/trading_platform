package data

import (
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
}

func NewServer() *Server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return &Server{
		Logger:     logger,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
		CacheTTL:   15 * time.Minute, // Set cache TTL to 15 minutes
	}
}

func (s *Server) getCachedData(cacheKey string) (*pb.StockResponse, bool) {
	if entry, ok := s.Cache.Load(cacheKey); ok {
		cacheEntry := entry.(CacheEntry)
		if time.Now().Before(cacheEntry.Expiration) {
			return cacheEntry.Data, true
		}
		// Cache entry has expired, remove it
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
