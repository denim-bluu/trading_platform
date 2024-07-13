// cmd/strategy_client/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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
	client := pb.NewStrategyServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter date (YYYY-MM-DD): ")
		date, _ := reader.ReadString('\n')
		date = strings.TrimSpace(date)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		r, err := client.GetTradingSignals(ctx, &pb.SignalRequest{Date: date})
		cancel()
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		fmt.Printf("Trading Signals for %s:\n", date)
		fmt.Printf("Market Regime: %t\n", r.IsMarketRegimePositive)
		for _, signal := range r.Signals {
			fmt.Printf("Stock: %s, Signal: %s, Position Size: %.2f\n",
				signal.Symbol, signal.Signal, signal.PositionSize)
		}
	}
}
