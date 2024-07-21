package strategy

import (
	"context"
	"fmt"
	pb "momentum-trading-platform/api/proto/strategy_service"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedStrategyServiceServer
	Logger     *log.Logger
	Clients    *Clients
	Strategies map[string]Strategy
}

func NewServer(clients *Clients) *Server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	s := &Server{
		Logger:     logger,
		Clients:    clients,
		Strategies: make(map[string]Strategy),
	}

	// Register the momentum strategy
	s.Strategies["momentum"] = NewMomentumStrategy()

	return s
}

func (s *Server) ConfigureStrategy(ctx context.Context, req *pb.ConfigureStrategyRequest) (*pb.ConfigureStrategyResponse, error) {
	params := make(map[string]interface{})
	for k, v := range req.Parameters {
		params[k] = v
	}

	err := s.configureStrategy(req.StrategyName, params)
	if err != nil {
		return &pb.ConfigureStrategyResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ConfigureStrategyResponse{
		Success: true,
		Message: "Strategy configured successfully",
	}, nil
}

func (s *Server) GetStrategyParameters(ctx context.Context, req *pb.GetStrategyParametersRequest) (*pb.GetStrategyParametersResponse, error) {
	params, err := s.getStrategyParameters(req.StrategyName)
	if err != nil {
		return nil, err
	}

	stringParams := make(map[string]string)
	for k, v := range params {
		stringParams[k] = fmt.Sprintf("%v", v)
	}

	return &pb.GetStrategyParametersResponse{
		Parameters: stringParams,
	}, nil
}

func (s *Server) configureStrategy(strategyName string, params map[string]interface{}) error {
	strategy, ok := s.Strategies[strategyName]
	if !ok {
		return fmt.Errorf("strategy %s not found", strategyName)
	}

	err := strategy.SetParameters(params)
	if err != nil {
		return fmt.Errorf("failed to set parameters for strategy %s: %v", strategyName, err)
	}

	s.Logger.WithFields(log.Fields{
		"strategy": strategyName,
		"params":   params,
	}).Info("Strategy configured")

	return nil
}

func (s *Server) getStrategyParameters(strategyName string) (map[string]interface{}, error) {
	strategy, ok := s.Strategies[strategyName]
	if !ok {
		return nil, fmt.Errorf("strategy %s not found", strategyName)
	}

	return strategy.GetParameters(), nil
}
