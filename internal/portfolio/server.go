// internal/portfolio/server.go
package portfolio

import (
	"context"
	"sync"
	"time"

	pb "momentum-trading-platform/api/proto/portfolio_service"
	portfoliostatepb "momentum-trading-platform/api/proto/portfolio_state_service"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedPortfolioServiceServer
	Logger            *log.Logger
	Clients           *Clients
	DesiredPortfolio  map[string]*pb.Position
	CashBalance       float64
	RebalanceSchedule string
	lastRebalanceTime time.Time
	mu                sync.Mutex
}

func NewServer(clients *Clients) (*Server, error) {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	s := &Server{
		Logger:            logger,
		Clients:           clients,
		DesiredPortfolio:  make(map[string]*pb.Position),
		RebalanceSchedule: "weekly",
		lastRebalanceTime: time.Now(),
	}

	// Load the previous state
	err := s.loadState()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) loadState() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the current portfolio state from the Portfolio State Service
	portfolioState, err := s.Clients.PortfolioStateClient.GetPortfolioState(ctx, &portfoliostatepb.GetPortfolioStateRequest{})
	if err != nil {
		s.Logger.WithError(err).Error("Failed to load portfolio state")
		return err
	}

	// Update the server's desired portfolio state based on the actual state
	s.DesiredPortfolio = make(map[string]*pb.Position)
	for _, position := range portfolioState.Positions {
		s.DesiredPortfolio[position.Symbol] = &pb.Position{
			Symbol:       position.Symbol,
			Quantity:     position.Quantity,
			CurrentPrice: position.CurrentPrice,
			MarketValue:  position.MarketValue,
		}
	}

	s.Logger.Info("Loaded current portfolio state as initial desired state")
	return nil
}

func (s *Server) calculateTotalValue() float64 {
	total := s.CashBalance
	for _, pos := range s.DesiredPortfolio {
		total += pos.MarketValue
	}
	return total
}
