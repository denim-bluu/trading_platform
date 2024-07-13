// cmd/backtesting_client/main.go
package main

import (
	"context"
	"time"

	"github.com/charmbracelet/log"

	pb "momentum-trading-platform/api/proto/backtesting_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewBacktestingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start a backtest
	backtestReq := &pb.BacktestRequest{
		StartDate:      "2022-01-01",
		EndDate:        "2024-12-31",
		InitialCapital: 100000,
		Symbols:        []string{"AAPL", "GOOGL", "MSFT", "AMZN", "FB"},
	}
	log.Infof("Starting backtest: %+v\n", backtestReq)
	result, err := client.RunBacktest(ctx, backtestReq)
	if err != nil {
		log.Fatalf("could not run backtest: %v", err)
	}
	log.Infof("Backtest started: %s\n", result.BacktestId)

	// Check status periodically
	for {
		status, err := client.GetBacktestStatus(ctx, &pb.BacktestStatusRequest{BacktestId: result.BacktestId})
		if err != nil {
			log.Fatalf("could not get backtest status: %v", err)
		}
		log.Infof("Backtest status: %s, Progress: %.2f%%\n", status.Status, status.Progress)

		if status.Status == "COMPLETED" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Get final result
	finalResult, err := client.GetBacktestStatus(ctx, &pb.BacktestStatusRequest{BacktestId: result.BacktestId})
	if err != nil {
		log.Fatalf("could not get final backtest result: %v", err)
	}
	log.Infof("Final backtest result: %+v\n", finalResult)
}
