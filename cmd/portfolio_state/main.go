// cmd/portfolio_state/main.go
package main

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "momentum-trading-platform/api/proto/portfolio_state_service"
	portfoliostate "momentum-trading-platform/internal/portfolio_state"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Infof("Connecting to database: %s", dbURI)

	s, err := portfoliostate.NewServer(dbURI)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPortfolioStateServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	s.Logger.WithField("address", lis.Addr().String()).Info("Portfolio state service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.WithError(err).Fatal("Failed to serve")
	}
}
