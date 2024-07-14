// cmd/portfolio_state/main.go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "momentum-trading-platform/api/proto/portfolio_state_service"
)

func main() {
	conn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPortfolioStateServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Example: Save a portfolio state
	saveResponse, err := savePortfolioState(ctx, c)
	if err != nil {
		log.Fatalf("could not save portfolio state: %v", err)
	}
	fmt.Printf("Save response: %s\n", saveResponse.Message)

	// Example: Load a portfolio state
	loadedState, err := loadPortfolioState(ctx, c, time.Now().Format("2006-01-02"))
	if err != nil {
		log.Fatalf("could not load portfolio state: %v", err)
	}
	fmt.Printf("Loaded portfolio state: %+v\n", loadedState)

	// Example: Get portfolio history
	history, err := getPortfolioHistory(ctx, c, time.Now().AddDate(0, -1, 0), time.Now())
	if err != nil {
		log.Fatalf("could not get portfolio history: %v", err)
	}
	fmt.Printf("Portfolio history: %+v\n", history)
}

func savePortfolioState(ctx context.Context, c pb.PortfolioStateServiceClient) (*pb.SaveResponse, error) {
	state := &pb.PortfolioState{
		Date: time.Now().Format("2006-01-02"),
		Positions: []*pb.Position{
			{
				Symbol:       "AAPL",
				Quantity:     100,
				AveragePrice: 150.0,
				CurrentPrice: 155.0,
				MarketValue:  15500.0,
			},
			{
				Symbol:       "GOOGL",
				Quantity:     50,
				AveragePrice: 2000.0,
				CurrentPrice: 2100.0,
				MarketValue:  105000.0,
			},
		},
		CashBalance: 50000.0,
		TotalValue:  170500.0,
	}

	return c.SavePortfolioState(ctx, state)
}

func loadPortfolioState(ctx context.Context, c pb.PortfolioStateServiceClient, date string) (*pb.PortfolioState, error) {
	req := &pb.LoadRequest{
		Date: date,
	}
	return c.LoadPortfolioState(ctx, req)
}

func getPortfolioHistory(ctx context.Context, c pb.PortfolioStateServiceClient, startDate, endDate time.Time) (*pb.PortfolioHistory, error) {
	req := &pb.HistoryRequest{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}
	return c.GetPortfolioHistory(ctx, req)
}
