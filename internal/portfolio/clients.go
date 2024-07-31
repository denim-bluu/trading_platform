// internal/portfolio/clients.go
package portfolio

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	datapb "momentum-trading-platform/api/proto/data_service"
	portfoliostatepb "momentum-trading-platform/api/proto/portfolio_state_service"
	strategypb "momentum-trading-platform/api/proto/strategy_service"
	tradepb "momentum-trading-platform/api/proto/trade_execution_service"
)

type Clients struct {
	DataClient           datapb.DataServiceClient
	StrategyClient       strategypb.StrategyServiceClient
	PortfolioStateClient portfoliostatepb.PortfolioStateServiceClient
	TradeExecutionClient tradepb.TradeExecutionServiceClient
	connections          []*grpc.ClientConn
}

func NewClients() (*Clients, error) {
	dataConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to data service: %v", err)
	}

	strategyConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to strategy service: %v", err)
	}

	portfolioStateConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to portfolio state service: %v", err)
	}

	return &Clients{
		DataClient:           datapb.NewDataServiceClient(dataConn),
		StrategyClient:       strategypb.NewStrategyServiceClient(strategyConn),
		PortfolioStateClient: portfoliostatepb.NewPortfolioStateServiceClient(portfolioStateConn),
		connections:          []*grpc.ClientConn{dataConn, strategyConn, portfolioStateConn},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.connections {
		conn.Close()
	}
}
