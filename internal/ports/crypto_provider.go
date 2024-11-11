package ports

import (
	"context"
	"kursy-kriptovalyut/internal/entity"
)

type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]entity.Coin, error)
}
