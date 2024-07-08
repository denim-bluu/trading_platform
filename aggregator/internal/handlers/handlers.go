package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"trading_platform/aggregator/internal/aggregator"
	pb "trading_platform/aggregator/proto"
)

type Handlers struct {
	Aggregator *aggregator.Aggregator
}

func NewHandlers(agg *aggregator.Aggregator) *Handlers {
	return &Handlers{Aggregator: agg}
}

func (h *Handlers) AggregateHistoricalDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to aggregate historical data")
	type request struct {
		Symbol    string `json:"symbol"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Aggregating historical data for symbol: %s, start date: %s, end date: %s", req.Symbol, req.StartDate, req.EndDate)
	err := h.Aggregator.AggregateHistoricalData(req.Symbol, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Failed to aggregate historical data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully aggregated historical data")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handlers) UpdateLiveDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to update live data")
	type request struct {
		Symbol string `json:"symbol"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updating live data for symbol: %s", req.Symbol)
	err := h.Aggregator.UpdateLiveData(req.Symbol)
	if err != nil {
		log.Printf("Failed to update live data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully updated live data")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handlers) GetHistoricalDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get historical data")
	type request struct {
		Symbol    string `json:"symbol"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Getting historical data for symbol: %s, start date: %s, end date: %s", req.Symbol, req.StartDate, req.EndDate)
	dataPoints, err := h.Aggregator.GetHistoricalData(req.Symbol, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Failed to get historical data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	historicalData := pb.HistoricalData{DataPoints: dataPoints}
	log.Println("Successfully retrieved historical data")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&historicalData)
}
