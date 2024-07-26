package tradeexecution

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	portfoliostatepb "momentum-trading-platform/api/proto/portfolio_state_service"
)

type Clients struct {
	PortfolioStateClient portfoliostatepb.PortfolioStateServiceClient
	connections          []*grpc.ClientConn
}

func NewClients() (*Clients, error) {
	portfolioStateConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to portfolio state service: %v", err)
	}

	return &Clients{
		PortfolioStateClient: portfoliostatepb.NewPortfolioStateServiceClient(portfolioStateConn),
		connections:          []*grpc.ClientConn{portfolioStateConn},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.connections {
		conn.Close()
	}
}
