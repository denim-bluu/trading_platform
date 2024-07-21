// cmd/data_client/main.go
package main

import (
	"context"
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
	c := pb.NewDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	symbol := "TXN"
	startDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	interval := "1d"

	log.Infof("Getting stock data for %s from %s to %s with interval %s", symbol, startDate, endDate, interval)
	r, err := c.GetStockData(ctx, &pb.StockRequest{Symbol: symbol, StartDate: startDate, EndDate: endDate, Interval: interval})
	if err != nil {
		log.Fatalf("could not get stock data: %v", err)
	}

	log.Infof("Total number of data points: %d", len(r.DataPoints))
	log.Printf("Stock Data for %s:", r.Symbol)
	for _, dp := range r.DataPoints {
		log.Printf("Time: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Adj Close: %.2f, Volume: %d",
			time.Unix(dp.Timestamp, 0).Format("2006-01-02"),
			dp.Open, dp.High, dp.Low, dp.Close, dp.AdjustedClose, dp.Volume)
	}

	// rs, err := c.GetBatchStockData(ctx, &pb.BatchStockRequest{Symbols: []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "NVDA", "NFLX", "PYPL", "ADBE",
	// 	"INTC", "CSCO", "CMCSA", "PEP", "AVGO", "TXN", "COST", "QCOM", "TMUS", "AMGN", "SBUX",
	// 	"INTU", "AMD", "ISRG", "GILD", "MDLZ", "BKNG", "MU", "ADP", "REGN", "ATVI"}, StartDate: startDate, EndDate: endDate, Interval: interval})
	// if err != nil {
	// 	log.Fatalf("could not get batch stock data: %v", err)
	// }

	// for symbol, resp := range rs.StockData {
	// 	log.Printf("Stock Data for %s:", symbol)
	// 	for _, dp := range resp.DataPoints {
	// 		log.Printf("Time: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Adj Close: %.2f, Volume: %d",
	// 			time.Unix(dp.Timestamp, 0).Format("2006-01-02"),
	// 			dp.Open, dp.High, dp.Low, dp.Close, dp.AdjustedClose, dp.Volume)
	// 	}
	// }
}
