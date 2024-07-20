package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/log"

	"momentum-trading-platform/internal/data"

	pb "momentum-trading-platform/api/proto/data_service"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Infof("Connecting to database: %s", dbURI)

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	s, err := data.NewServer(db)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		s.Logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, s)

	s.Logger.WithField("address", lis.Addr().String()).Info("Data service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.WithError(err).Fatal("Failed to serve")
	}
}
