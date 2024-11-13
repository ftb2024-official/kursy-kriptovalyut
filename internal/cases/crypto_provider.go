package cases

import (
	"context"
	entity "kursy-kriptovalyut/internal/entities"
)

type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]entity.Coin, error)
}

// GetActualRates - срабатывает, когда запрашиваемый коин отсутсвует в БД
