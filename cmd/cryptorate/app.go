package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"kursy-kriptovalyut/internal/adapters/provider"
	"kursy-kriptovalyut/internal/adapters/storage"
	"kursy-kriptovalyut/internal/cases"
	"kursy-kriptovalyut/internal/ports"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	provider, err := provider.NewCryptoCompare("https://min-api.cryptocompare.com/data/pricemultifull", "851e396ad68e892830b474f074b051d2104b77576c25b9058ef16d4a477515d8")
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewPostgres("postgres://user:pswd@localhost:5434/crypto_rate?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	service, err := cases.NewService(provider, storage)
	if err != nil {
		log.Fatal(err)
	}

	server, err := ports.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	go startTicker(ctx, service)

	log.Println("Server running on port :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func startTicker(ctx context.Context, service *cases.Service) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// fmt.Println("tick-tack")
			if err := service.ActualizeRates(ctx); err != nil {
				log.Println("Failed to actualize rates:", err)
			}
		case <-ctx.Done():
			log.Println("Stopping ticker...")
			return
		}
	}
}
