// cmd/portfolio_client/main.go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	pb "momentum-trading-platform/api/proto/portfolio_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewPortfolioServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Get initial portfolio status
	status, err := client.GetPortfolioStatus(ctx, &pb.PortfolioRequest{Date: time.Now().Format("2006-01-02")})
	if err != nil {
		log.Fatalf("could not get portfolio status: %v", err)
	}
	fmt.Printf("Initial Portfolio Status: %+v\n", status)

	// Process some mock trading signals
	update, err := client.ProcessTradingSignals(ctx, &pb.TradingSignals{
		Signals: []*pb.Signal{
			{Symbol: "AAPL", Type: pb.TradeType_BUY, PositionSize: 1},
			{Symbol: "GOOGL", Type: pb.TradeType_BUY, PositionSize: 5},
		},
	})
	if err != nil {
		log.Fatalf("could not process trading signals: %v", err)
	}
	fmt.Printf("Portfolio Update After Signals: %+v\n", update)

	// Perform a rebalance
	rebalanceUpdate, err := client.RebalancePortfolio(ctx, &pb.RebalanceRequest{Date: time.Now().Format("2006-01-02")})
	if err != nil {
		log.Fatalf("could not rebalance portfolio: %v", err)
	}
	fmt.Printf("Portfolio Update After Rebalance: %+v\n", rebalanceUpdate)
}
