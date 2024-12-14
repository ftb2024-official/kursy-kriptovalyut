package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"kursy-kriptovalyut/internal/adapters/provider"
	"kursy-kriptovalyut/internal/adapters/storage"
	"kursy-kriptovalyut/internal/app"
	"kursy-kriptovalyut/internal/cases"
	"kursy-kriptovalyut/internal/ports"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := app.NewEnv()
	baseUrl := env.Url
	apiKey := env.ApiKey
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", env.PgUser, env.PgPswd, env.PgHost, env.PgPort, env.PgDB)

	provider, err := provider.NewCryptoCompare(baseUrl, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewPostgres(connStr)
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
		Addr:    env.SrvPort,
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

	log.Printf("Server running on port %v", env.SrvPort)
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
			fmt.Println("tick-tack")
			if err := service.ActualizeRates(ctx); err != nil {
				log.Println("Failed to actualize rates:", err)
			}
		case <-ctx.Done():
			log.Println("Stopping ticker...")
			return
		}
	}
}
