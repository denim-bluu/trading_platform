// cmd/strategy_client/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "momentum-trading-platform/api/proto/strategy_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStrategyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "NVDA", "PYPL", "INTC", "NFLX", "ADBE"}
	startDate := fmt.Sprintf("%d", time.Now().AddDate(-1, 0, 0).Unix())
	endDate := fmt.Sprintf("%d", time.Now().Unix())
	interval := "1d"

	req := &pb.SignalRequest{
		Symbols:   symbols,
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  interval,
	}

	resp, err := c.GenerateSignals(ctx, req)
	if err != nil {
		log.Fatalf("could not generate signals: %v", err)
	}

	log.Println("Generated Signals:")
	for _, signal := range resp.Signals {
		log.Printf("Symbol: %s, Signal: %s, Position Size: %.2f, Momentum Score: %.4f",
			signal.Symbol, signal.Signal, signal.PositionSize, signal.MomentumScore)
	}
}
