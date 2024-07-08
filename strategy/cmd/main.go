package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trading_platform/strategy/internal/config"
	"trading_platform/strategy/internal/data"
	"trading_platform/strategy/internal/handlers"
	"trading_platform/strategy/middleware"
)

func main() {
	log.Println("Starting the Strategy Service...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fetcher := data.NewHTTPDataFetcher(cfg.AGGREGATOR_URL)
	handler := handlers.NewHandlers(fetcher)

	mux := http.NewServeMux()
	mux.HandleFunc("/evaluate_strategy", middleware.Chain(
		handler.EvaluateStrategyHandler,
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
