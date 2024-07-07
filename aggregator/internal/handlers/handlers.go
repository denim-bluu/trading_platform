package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"trading_platform/aggregator/internal/aggregator"
)

type Handlers struct {
	Aggregator *aggregator.Aggregator
}

func NewHandlers(agg *aggregator.Aggregator) *Handlers {
	return &Handlers{Aggregator: agg}
}
func (h *Handlers) AggregateHistoricalDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Symbol    string `json:"symbol"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Filename  string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Symbol == "" || req.StartDate == "" || req.EndDate == "" || req.Filename == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	log.Printf("Aggregating historical data for symbol: %s, start date: %s, end date: %s, filename: %s", req.Symbol, req.StartDate, req.EndDate, req.Filename)
	err := h.Aggregator.AggregateHistoricalData(req.Symbol, req.StartDate, req.EndDate, req.Filename)
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
		Symbol   string `json:"symbol"`
		Filename string `json:"filename"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updating live data for symbol: %s, filename: %s", req.Symbol, req.Filename)
	err := h.Aggregator.UpdateLiveData(req.Symbol, req.Filename)
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
		Filename  string `json:"filename"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Getting historical data for symbol: %s, filename: %s, start date: %s, end date: %s", req.Symbol, req.Filename, req.StartDate, req.EndDate)
	dataPoints, err := h.Aggregator.GetHistoricalData(req.Symbol, req.Filename, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Failed to get historical data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully retrieved historical data")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataPoints)
}
