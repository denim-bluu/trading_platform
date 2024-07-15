package strategy

import (
	pb "momentum-trading-platform/api/proto/strategy_service"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedStrategyServiceServer
	Logger  *log.Logger
	Clients *Clients
}

func NewServer(clients *Clients) *Server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	return &Server{
		Logger:  logger,
		Clients: clients,
	}
}
