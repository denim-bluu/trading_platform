// cmd/backtesting/main.go
package main

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/charmbracelet/log"

	pb "momentum-trading-platform/api/proto/backtesting_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedBacktestingServiceServer
	backtests map[string]*backtestJob
	mu        sync.Mutex
}

type backtestJob struct {
	status   string
	progress float64
	result   *pb.BacktestResult
}

func newServer() *server {
	return &server{
		backtests: make(map[string]*backtestJob),
	}
}

func (s *server) RunBacktest(ctx context.Context, req *pb.BacktestRequest) (*pb.BacktestResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	backtestID := generateBacktestID()
	job := &backtestJob{
		status:   "RUNNING",
		progress: 0,
		result:   &pb.BacktestResult{BacktestId: backtestID},
	}
	s.backtests[backtestID] = job
	log.Debugf("Running backtest %s", backtestID)

	go s.runBacktestSimulation(backtestID, req)

	return &pb.BacktestResult{BacktestId: backtestID, Status: &pb.BacktestStatus{BacktestId: backtestID, Status: "RUNNING", Progress: 0}}, nil
}

func (s *server) GetBacktestStatus(ctx context.Context, req *pb.BacktestStatusRequest) (*pb.BacktestStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.backtests[req.BacktestId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "backtest not found")
	}

	return &pb.BacktestStatus{
		BacktestId: req.BacktestId,
		Status:     job.status,
		Progress:   job.progress,
	}, nil
}

func (s *server) runBacktestSimulation(backtestID string, req *pb.BacktestRequest) {
	log.Infof("Running backtest simulation for %s", backtestID)
	// Simulate backtesting process
	time.Sleep(5 * time.Second)

	s.mu.Lock()
	defer s.mu.Unlock()

	job := s.backtests[backtestID]
	job.status = "COMPLETED"
	job.progress = 100

	// Generate mock backtest results
	job.result = &pb.BacktestResult{
		BacktestId:          backtestID,
		Status:              &pb.BacktestStatus{BacktestId: backtestID, Status: "COMPLETED", Progress: 100},
		FinalPortfolioValue: req.InitialCapital * (1 + rand.Float64()),
		TotalReturn:         rand.Float64() * 0.5,
		SharpeRatio:         1 + rand.Float64(),
		MaxDrawdown:         rand.Float64() * 0.2,
		Trades:              generateMockTrades(req),
	}
}

func generateMockTrades(req *pb.BacktestRequest) []*pb.TradeRecord {
	trades := make([]*pb.TradeRecord, 0)
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		if rand.Float32() < 0.1 {
			tradeRecord := &pb.TradeRecord{
				Date:     date.Format("2006-01-02"),
				Symbol:   req.Symbols[rand.Intn(len(req.Symbols))],
				Action:   []string{"BUY", "SELL"}[rand.Intn(2)],
				Quantity: rand.Int31n(100) + 1,
				Price:    100 + rand.Float64()*100,
			}
			log.Infof("Generated trade: %+v", tradeRecord)
			trades = append(trades, tradeRecord)
		}
	}

	return trades
}

func generateBacktestID() string {
	return time.Now().Format("20060102150405")
}

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBacktestingServiceServer(s, newServer())
	log.Printf("Backtesting service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
