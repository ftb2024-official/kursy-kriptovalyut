// package main

// import (
// 	"context"
// 	"fmt"
// 	cfg "kursy-kriptovalyut"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"time"

// 	"kursy-kriptovalyut/internal/adapters/provider"
// 	"kursy-kriptovalyut/internal/adapters/storage"
// 	"kursy-kriptovalyut/internal/cases"
// 	"kursy-kriptovalyut/internal/ports"
// 	"kursy-kriptovalyut/pkg/logger"

// 	"go.uber.org/zap"
// )

// var LOG = logger.NewLogger()

// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	cfg := cfg.LoadCfg()

// 	fmt.Println(cfg)

// 	baseUrl := cfg.Cfg.Url
// 	apiKey := cfg.Cfg.ApiKey
// 	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Cfg.PgUser, cfg.Cfg.PgPswd, cfg.Cfg.PgHost, cfg.Cfg.PgPort, cfg.Cfg.PgDB)

// 	provider, err := provider.NewCryptoCompare(baseUrl, apiKey)
// 	if err != nil {
// 		LOG.Fatal("failed to create provider", zap.Any("err", err.Error()))
// 	}

// 	storage, err := storage.NewPostgres(connStr)
// 	if err != nil {
// 		LOG.Fatal("failed to create storage", zap.Any("err", err.Error()))
// 	}

// 	service, err := cases.NewService(provider, storage)
// 	if err != nil {
// 		LOG.Fatal("failed to create service", zap.Any("err", err.Error()))
// 	}

// 	server, err := ports.NewServer(service)
// 	if err != nil {
// 		LOG.Fatal("failed to create server", zap.Any("err", err.Error()))
// 	}

// 	srv := &http.Server{
// 		Addr:    cfg.Cfg.Port,
// 		Handler: server,
// 	}

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, os.Interrupt)

// 	go func() {
// 		<-quit
// 		LOG.Info("Shutting down server...")

// 		ctx, cancel := context.WithCancel(context.Background())
// 		defer cancel()

// 		if err := srv.Shutdown(ctx); err != nil {
// 			LOG.Error("server shutdown", zap.Any("err", err.Error()))
// 		}
// 	}()

// 	go startTicker(ctx, service)

// 	LOG.Info("Server running on port :8080")
// 	if err := srv.ListenAndServe(); err != nil || err != http.ErrServerClosed {
// 		LOG.Error("failed to start server", zap.Any("err", err.Error()))
// 	}
// }

// func startTicker(ctx context.Context, service *cases.Service) {
// 	ticker := time.NewTicker(time.Minute * 1)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			LOG.Info("TICK-TACK")
// 			if err := service.ActualizeRates(ctx); err != nil {
// 				LOG.Info("failed to actualize rates:", zap.Any("err", err.Error()))
// 			}
// 		case <-ctx.Done():
// 			LOG.Info("stopping ticker...")
// 			return
// 		}
// 	}
// }

package main

import "kursy-kriptovalyut/internal/app"

func main() {
	app := app.NewApp()
	app.Run()

}
