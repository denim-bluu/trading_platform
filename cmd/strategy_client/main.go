// cmd/strategy_client/main.go
package main

import (
	"context"
	"time"

	"github.com/charmbracelet/log"

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

	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "NVDA", "NFLX", "PYPL", "ADBE",
		"INTC", "CSCO", "CMCSA", "PEP", "AVGO", "TXN", "COST", "QCOM", "TMUS", "AMGN", "SBUX",
		"INTU", "AMD", "ISRG", "GILD", "MDLZ", "BKNG", "MU", "ADP", "REGN", "ATVI"}
	startDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	interval := "1d"

	req := &pb.SignalRequest{
		Symbols:     symbols,
		StartDate:   startDate,
		EndDate:     endDate,
		Interval:    interval,
		MarketIndex: "^GSPC",
	}

	resp, err := c.GenerateSignals(ctx, req)
	if err != nil {
		log.Fatalf("could not generate signals: %v", err)
	}

	log.Infof("Generated %d signals\n", len(resp.Signals))
	log.Print("Generated Signals:")
	for _, signal := range resp.Signals {
		log.Infof("Symbol: %s, Signal: %s, Risk Unit: %.4f, Momentum Score: %.4f, Last Close: %.2f\n",
			signal.Symbol, signal.Signal, signal.RiskUnit, signal.MomentumScore, signal.CurrentPrice)
	}
}
