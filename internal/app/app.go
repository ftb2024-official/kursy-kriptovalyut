package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/robfig/cron"
	"go.uber.org/zap"

	"kursy-kriptovalyut/config"
	"kursy-kriptovalyut/internal/adapters/provider"
	"kursy-kriptovalyut/internal/adapters/storage"
	"kursy-kriptovalyut/internal/cases"
	"kursy-kriptovalyut/internal/ports"
	"kursy-kriptovalyut/pkg/logger"
)

var log = logger.NewLogger()

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	cfg := config.LoadCfg()
	baseUrl := cfg.Cfg.Url
	apiKey := cfg.Cfg.ApiKey
	user := cfg.Cfg.PgUser
	pswd := cfg.Cfg.PgPswd
	host := cfg.Cfg.PgHost
	port := cfg.Cfg.PgPort
	db := cfg.Cfg.PgDB
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, pswd, host, port, db)

	provider, err := provider.NewCryptoCompare(baseUrl, apiKey)
	if err != nil {
		log.Fatal("failed to create provider", zap.Any("err", err.Error()))
	}

	storage, err := storage.NewPostgres(connStr)
	if err != nil {
		log.Fatal("failed to create storage", zap.Any("err", err.Error()))
	}

	service, err := cases.NewService(provider, storage)
	if err != nil {
		log.Fatal("failed to create service", zap.Any("err", err.Error()))
	}

	ctx := context.Background()
	cron := SetCron(service, ctx)
	cron.Start()

	server, err := ports.NewServer(service)
	if err != nil {
		log.Fatal("failed to create server", zap.Any("err", err.Error()))
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Info("Shutting down server...")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server shutdown", zap.Any("err", err.Error()))
		}
	}()

	log.Info("Server running on port :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("failed to start server", zap.Any("err", err.Error()))
	}
}

func SetCron(srvc *cases.Service, ctx context.Context) *cron.Cron {
	c := cron.New()
	err := c.AddFunc("@every 1m", func() { _ = srvc.ActualizeRates(ctx) })
	if err != nil {
		log.Error("cron job failed", zap.Any("err", err.Error()))
	}

	return c
}
