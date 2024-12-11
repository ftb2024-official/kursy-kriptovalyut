package ports

import (
	"context"

	"kursy-kriptovalyut/internal/entities"
)

type Service interface {
	GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAggRates(ctx context.Context, titles []string, aggFuncName string) ([]entities.Coin, error)
}
