// internal/portfolio/server.go

package portfolio

import (
	"context"
	"os"
	"sync"
	"time"

	pb "momentum-trading-platform/api/proto/portfolio_service"
	"momentum-trading-platform/internal/storage"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedPortfolioServiceServer
	Logger         *log.Logger
	Clients        *Clients
	Portfolio      map[string]*pb.Position
	Storage        storage.Storage
	CashBalance    float64
	LastUpdateDate string
	mu             sync.Mutex
}

func NewServer(clients *Clients, storage storage.Storage) *Server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	logger.SetOutput(os.Stdout)

	s := &Server{
		Logger:         logger,
		Clients:        clients,
		Portfolio:      make(map[string]*pb.Position),
		Storage:        storage,
		LastUpdateDate: time.Now().Format("2006-01-02:15:04:05"),
		CashBalance:    1000000, // Starting with $1M cash
	}

	// Load initial state from storage
	if state, err := storage.LoadPortfolioState(context.Background()); err == nil {
		s.Portfolio = make(map[string]*pb.Position)
		for _, position := range state.Positions {
			s.Portfolio[position.Symbol] = position
		}
		s.CashBalance = state.CashBalance
		s.LastUpdateDate = state.LastUpdateDate
	}

	return s
}
