// cmd/portfolio_client/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "momentum-trading-platform/api/proto/portfolio_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPortfolioServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Example: Perform weekly rebalance
	weeklyUpdate, err := performWeeklyRebalance(ctx, c)
	if err != nil {
		log.Fatalf("could not perform weekly rebalance: %v", err)
	}
	fmt.Printf("Weekly rebalance result: %+v\n", weeklyUpdate)

	// Example: Perform bi-weekly rebalance
	biWeeklyUpdate, err := performBiWeeklyRebalance(ctx, c)
	if err != nil {
		log.Fatalf("could not perform bi-weekly rebalance: %v", err)
	}
	fmt.Printf("Bi-weekly rebalance result: %+v\n", biWeeklyUpdate)

	// Example: Get portfolio status
	status, err := getPortfolioStatus(ctx, c)
	if err != nil {
		log.Fatalf("could not get portfolio status: %v", err)
	}
	fmt.Printf("Portfolio status: %+v\n", status)
}

func performWeeklyRebalance(ctx context.Context, c pb.PortfolioServiceClient) (*pb.PortfolioUpdate, error) {
	req := &pb.RebalanceRequest{
		Date: time.Now().Format("2006-01-02"),
	}
	return c.WeeklyRebalance(ctx, req)
}

func performBiWeeklyRebalance(ctx context.Context, c pb.PortfolioServiceClient) (*pb.PortfolioUpdate, error) {
	req := &pb.RebalanceRequest{
		Date: time.Now().Format("2006-01-02"),
	}
	return c.BiWeeklyRebalance(ctx, req)
}

func getPortfolioStatus(ctx context.Context, c pb.PortfolioServiceClient) (*pb.PortfolioStatus, error) {
	req := &pb.PortfolioStatusRequest{}
	return c.GetPortfolioStatus(ctx, req)
}
