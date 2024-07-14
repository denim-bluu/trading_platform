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
	c := pb.NewPortfolioServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Get current portfolio status
	status, err := c.GetPortfolioStatus(ctx, &pb.PortfolioStatusRequest{})
	if err != nil {
		log.Fatalf("could not get portfolio status: %v", err)
	}
	printPortfolioStatus(status)

	// Update portfolio
	updateReq := &pb.UpdatePortfolioRequest{
		Date: time.Now().Format("2006-01-02"),
	}
	log.Info("Updating portfolio...")
	updateResp, err := c.UpdatePortfolio(ctx, updateReq)
	if err != nil {
		log.Fatalf("could not update portfolio: %v", err)
	}
	printPortfolioUpdate(updateResp)

	// // Rebalance portfolio
	// rebalanceReq := &pb.RebalanceRequest{
	// 	Date: time.Now().Format("2006-01-02"),
	// }
	// rebalanceResp, err := c.RebalancePortfolio(ctx, rebalanceReq)
	// if err != nil {
	// 	log.Fatalf("could not rebalance portfolio: %v", err)
	// }
	// printPortfolioUpdate(rebalanceResp)
}

func printPortfolioStatus(status *pb.PortfolioStatus) {
	fmt.Println("Portfolio Status:")
	fmt.Printf("Cash Balance: $%.2f\n", status.CashBalance)
	fmt.Printf("Total Value: $%.2f\n", status.TotalValue)
	fmt.Printf("Last Update Date: %s\n", status.LastUpdateDate)
	fmt.Println("Positions:")
	for _, pos := range status.Positions {
		fmt.Printf("  %s: %d shares, Avg Price: $%.2f, Current Price: $%.2f, Market Value: $%.2f\n",
			pos.Symbol, pos.Quantity, pos.AveragePrice, pos.CurrentPrice, pos.MarketValue)
	}
	fmt.Println()
}

func printPortfolioUpdate(update *pb.PortfolioUpdate) {
	fmt.Println("Portfolio Update:")
	fmt.Println("Trades:")
	for _, trade := range update.Trades {
		fmt.Printf("  %s %s: %d shares at $%.2f\n",
			trade.Type, trade.Symbol, trade.Quantity, trade.Price)
	}
	fmt.Println("Updated Status:")
	printPortfolioStatus(update.UpdatedStatus)
}
