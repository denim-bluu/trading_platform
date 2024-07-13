// cmd/strategy/main.go
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"time"

	datapb "momentum-trading-platform/api/proto/data_service"
	pb "momentum-trading-platform/api/proto/strategy_service"
	"momentum-trading-platform/utils"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	pb.UnimplementedStrategyServiceServer
	logger     *log.Logger
	dataClient datapb.DataServiceClient
}

func parseAndFormatTimestamp(timestampStr string) (string, error) {
	// Parse int64 from string
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Format the timestamp into a date string
	date := time.Unix(timestamp, 0).Format("2006-01-02:15:04:05")
	return date, nil
}

func NewServer(dataClient datapb.DataServiceClient) *Server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	return &Server{
		logger:     logger,
		dataClient: dataClient,
	}
}

func (s *Server) GenerateSignals(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
	startDate, err := parseAndFormatTimestamp(req.StartDate)
	if err != nil {
		s.logger.WithError(err).Error("Failed to parse start date")
		return nil, fmt.Errorf("failed to parse start date: %v", err)
	}
	endDate, err := parseAndFormatTimestamp(req.EndDate)
	if err != nil {
		s.logger.WithError(err).Error("Failed to parse end date")
		return nil, fmt.Errorf("failed to parse end date: %v", err)
	}
	s.logger.WithFields(log.Fields{
		"symbols":  req.Symbols,
		"start":    startDate,
		"end":      endDate,
		"interval": req.Interval,
	}).Info("📧 Sending request to data service")

	batchReq := &datapb.BatchStockRequest{
		Symbols:   req.Symbols,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Interval:  req.Interval,
	}

	batchResp, err := s.dataClient.GetBatchStockData(ctx, batchReq)
	if err != nil {
		s.logger.WithError(err).Error("Failed to fetch batch stock data")
		return nil, fmt.Errorf("failed to fetch stock data: %v", err)
	}
	s.logger.Info("✅ Received batch stock data")

	signals := make([]*pb.StockSignal, 0, len(batchResp.StockData))
	for symbol, stockResp := range batchResp.StockData {
		s.logger.WithField("symbol", symbol).Info("📊 Processing stock data")

		// Disqualify stock if it has a recent large gap
		if utils.HasRecentLargeGap(stockResp.DataPoints, 90, 0.15) {
			s.logger.WithField("symbol", symbol).Info("❌ Stock disqualified due to recent large gap")
			continue
		}

		// Momentum score is for ranking stocks
		momentumScore := utils.CalculateMomentumScore(stockResp.DataPoints, 90)

		// ATR is for position sizing
		atr := utils.CalculateATR(stockResp.DataPoints, 20)
		positionSize := utils.CalculatePositionSize(atr, 1000000) // Assuming $1M account value

		// Latest Price and 100MA for signal generation
		lastPrice := stockResp.DataPoints[len(stockResp.DataPoints)-1].Close
		movingAverage := utils.CalculateMovingAverage(stockResp.DataPoints, 100)
		signal := utils.GenerateSignal(momentumScore, lastPrice, movingAverage)

		s.logger.WithField("momentum_score", momentumScore).Trace("🔢 Momentum score")
		s.logger.WithField("atr", atr).Trace("🔢 ATR")
		s.logger.WithField("position_size", positionSize).Trace("🔢 Position size")
		s.logger.WithField("last_price", lastPrice).Trace("🔢 Last price")
		s.logger.WithField("moving_average", movingAverage).Trace("🔢 Moving average")
		s.logger.WithField("signal", signal).Trace("🔢 Signal")

		signals = append(signals, &pb.StockSignal{
			Symbol:        symbol,
			Signal:        signal,
			PositionSize:  positionSize,
			MomentumScore: momentumScore,
		})
	}

	// Sort signals by momentum score in descending order
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].MomentumScore > signals[j].MomentumScore
	})
	s.logger.WithField("signals", signals).Info("Rank ordered Signals")

	// Keep only top 20% of stocks
	topCount := int(float64(len(signals)) * 0.2)
	signals = signals[:topCount]

	return &pb.SignalResponse{
		Signals: signals,
	}, nil
}

func main() {
	dataConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to data service: %v", err)
	}
	defer dataConn.Close()
	dataClient := datapb.NewDataServiceClient(dataConn)

	s := NewServer(dataClient)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		s.logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStrategyServiceServer(grpcServer, s)

	s.logger.WithField("address", lis.Addr().String()).Info("Strategy service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.WithError(err).Fatal("Failed to serve")
	}
}
