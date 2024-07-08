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

	log.Printf("Evaluating strategy for symbol: %s, strategy type: %s, start date: %s, end date: %s", req.Symbol, req.StrategyType, req.StartDate, req.EndDate)
	data, err := h.fetcher.FetchData(req.Symbol, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indexValue, err := h.fetcher.FetchIndexValue("^GSPC")
	if err != nil {
		log.Printf("Failed to fetch index value: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var strategy strategies.TradingStrategy
	switch req.StrategyType {
	case "momentum":
		strategy = strategies.NewMomentumStrategy()
	default:
		log.Printf("Unknown strategy type: %s", req.StrategyType)
		http.Error(w, "Unknown strategy type", http.StatusBadRequest)
		return
	}

	tradeActions := strategy.Evaluate(data.DataPoints, req.StartDate, req.EndDate, indexValue)
	res := &pb.StrategyResponse{TradeActions: tradeActions}

	log.Println("Successfully evaluated strategy")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
