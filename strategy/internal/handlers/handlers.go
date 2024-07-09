package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"trading_platform/strategy/internal/data"
	"trading_platform/strategy/internal/strategies"
	pb "trading_platform/strategy/proto"
)

type Handlers struct {
	fetcher data.DataFetcher
}

func NewHandlers(fetcher data.DataFetcher) *Handlers {
	return &Handlers{fetcher: fetcher}
}

func (h *Handlers) EvaluateStrategyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to evaluate strategy")
	var req pb.StrategyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch account value from the request or set a default value
	accountValue := req.AccountValue
	if accountValue <= 0 {
		accountValue = 100000.0 // Default account value
	}

	log.Printf("Evaluating strategy for symbol: %s, strategy type: %s, start date: %s, end date: %s, account value: %f", req.Symbol, req.StrategyType, req.StartDate, req.EndDate, accountValue)
	data, err := h.fetcher.FetchData(req.Symbol, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indexData, err := h.fetcher.FetchData("^GSPC", req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Failed to fetch index data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	strategy := strategies.NewMomentumStrategy()
	tradeActions := strategy.Evaluate(data.DataPoints, req.StartDate, req.EndDate, indexData.DataPoints, accountValue)
	res := &pb.StrategyResponse{TradeActions: tradeActions}

	log.Println("Successfully evaluated strategy")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
