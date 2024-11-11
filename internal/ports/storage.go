package ports

import (
	"context"
	"kursy-kriptovalyut/internal/entity"
)

type Storage interface {
	Store(ctx context.Context, coins []entity.Coin) error
	GetCoinsList(ctx context.Context) ([]string, error)
	GetActualCoins(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetAggregateCoins(ctx context.Context, titles []string) ([]entity.Coin, error)
}
