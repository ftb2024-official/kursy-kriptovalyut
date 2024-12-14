package app

import (
	"fmt"
	"log"
	"net/http"

	"kursy-kriptovalyut/internal/adapters/provider"
	"kursy-kriptovalyut/internal/adapters/storage"
	"kursy-kriptovalyut/internal/cases"
	"kursy-kriptovalyut/internal/ports"
)

type App struct {
	srv *ports.Server
}

func NewApp() *App {
	env := NewEnv()
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

	return &App{srv: server}
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: a.srv,
	}

	log.Printf("Server running on port :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
