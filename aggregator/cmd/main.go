package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trading_platform/aggregator/internal/aggregator"
	"trading_platform/aggregator/internal/config"
	"trading_platform/aggregator/internal/data"
	"trading_platform/aggregator/internal/handlers"
	"trading_platform/aggregator/internal/indicators"
	"trading_platform/aggregator/middleware"
)

func main() {
	log.Println("Starting the Aggregator Service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fetcher := data.NewYahooFinanceFetcher(cfg.YahooFinanceAPIURL)
	calculator := indicators.NewIndicatorCalculator()
	storage := data.NewDataStorage()

	agg := aggregator.NewAggregator(fetcher, calculator, storage)
	handler := handlers.NewHandlers(agg)

	mux := http.NewServeMux()
	mux.HandleFunc("/aggregate_historical_data", middleware.Chain(
		handler.AggregateHistoricalDataHandler,
		middleware.Logging(),
		middleware.Recovery(),
	))
	mux.HandleFunc("/update_live_data", middleware.Chain(
		handler.UpdateLiveDataHandler,
		middleware.Logging(),
		middleware.Recovery(),
	))
	mux.HandleFunc("/get_historical_data", middleware.Chain(
		handler.GetHistoricalDataHandler,
		middleware.Logging(),
		middleware.Recovery(),
	))

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
