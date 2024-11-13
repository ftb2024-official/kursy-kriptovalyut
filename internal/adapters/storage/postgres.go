package storage

import (
	"context"
	"database/sql"
	entity "kursy-kriptovalyut/internal/entities"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Store(ctx context.Context, coins []entity.Coin) error {
	return nil
}

func (p *Postgres) GetCoinsList(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (p *Postgres) GetActualCoins(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}

func (p *Postgres) GetAggregateCoins(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}
