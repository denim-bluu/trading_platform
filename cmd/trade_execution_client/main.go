// cmd/trade_execution_client/main.go
package main

import (
	"context"
	"time"

	"github.com/charmbracelet/log"

	pb "momentum-trading-platform/api/proto/trade_execution_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTradeExecutionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Execute a market order
	marketOrder := &pb.TradeOrder{
		Symbol:   "AAPL",
		Type:     pb.OrderType_MARKET,
		Quantity: 100,
		OrderId:  "order1",
	}
	result, err := client.ExecuteTrade(ctx, marketOrder)
	if err != nil {
		log.Fatalf("could not execute trade: %v", err)
	}
	log.Debugf("Market Order Result: %+v\n", result)

	// Execute a limit order
	limitOrder := &pb.TradeOrder{
		Symbol:     "GOOGL",
		Type:       pb.OrderType_LIMIT,
		Quantity:   200,
		LimitPrice: 150.0,
		OrderId:    "order2",
	}
	result, err = client.ExecuteTrade(ctx, limitOrder)
	if err != nil {
		log.Fatalf("could not execute trade: %v", err)
	}
	log.Debugf("Limit Order Result: %+v\n", result)

	// Check status of an order
	status, err := client.GetExecutionStatus(ctx, &pb.ExecutionStatusRequest{OrderId: "order1"})
	if err != nil {
		log.Fatalf("could not get execution status: %v", err)
	}
	log.Debugf("Order Status: %+v\n", status)
}
