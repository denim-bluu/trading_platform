package portfolio

import (
	"context"
	"sync"
	"time"

	pb "momentum-trading-platform/api/proto/portfolio_service"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedPortfolioServiceServer
	Logger         *log.Logger
	Clients        *Clients
	Portfolio      map[string]*pb.Position
	CashBalance    float64
	LastUpdateDate string
	mu             sync.Mutex
}

func NewServer(clients *Clients) *Server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	s := &Server{
		Logger:         logger,
		Clients:        clients,
		Portfolio:      make(map[string]*pb.Position),
		LastUpdateDate: time.Now().Format("2006-01-02:15:04:05"),
		CashBalance:    1000000, // Starting with $1M cash
	}

	if err := s.loadLastState(); err != nil {
		s.Logger.WithError(err).Error("Failed to load last portfolio state")
	}

	return s
}

func (s *Server) GetPortfolioStatus(ctx context.Context, req *pb.PortfolioStatusRequest) (*pb.PortfolioStatus, error) {
	return s.getPortfolioStatus(), nil
}

func (s *Server) getPortfolioStatus() *pb.PortfolioStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	positions := make([]*pb.Position, 0, len(s.Portfolio))
	for _, position := range s.Portfolio {
		positions = append(positions, position)
	}
	return &pb.PortfolioStatus{
		Positions:      positions,
		CashBalance:    s.CashBalance,
		TotalValue:     s.getTotalPortfolioValue(),
		LastUpdateDate: s.LastUpdateDate,
	}
}

func (s *Server) getTotalPortfolioValue() float64 {
	total := s.CashBalance
	for _, position := range s.Portfolio {
		total += position.MarketValue
	}
	return total
}
