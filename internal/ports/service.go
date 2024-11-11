package ports

import (
	"context"
	"kursy-kriptovalyut/internal/entity"
)

type Service interface {
	GetLastRates(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMaxRates(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMinRates(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetAvgRates(ctx context.Context, titles []string) ([]entity.Coin, error)
}
