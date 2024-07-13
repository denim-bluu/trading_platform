package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	pb "momentum-trading-platform/api/proto/data_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewDataServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter stock symbol (or 'quit' to exit): ")
		symbol, _ := reader.ReadString('\n')
		symbol = strings.TrimSpace(symbol)

		if symbol == "quit" {
			break
		}

		fmt.Print("Enter date (YYYY-MM-DD): ")
		date, _ := reader.ReadString('\n')
		date = strings.TrimSpace(date)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		r, err := client.GetStockData(ctx, &pb.StockRequest{Symbol: symbol, Date: date})
		cancel()
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		fmt.Printf("Stock Data: Symbol: %s, Date: %s, Price: %.2f, Volume: %.0f\n",
			r.Symbol, r.Date, r.Price, r.Volume)
	}
}
