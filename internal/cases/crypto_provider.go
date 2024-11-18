package cases

import (
	"context"
	entity "kursy-kriptovalyut/internal/entities"
)

//go:generate mockgen -source=./crypto_provider.go -destination=./mocks/gen/mock_crypto_provider.go
type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]entity.Coin, error)
}

// GetActualRates - срабатывает, когда запрашиваемый коин отсутсвует в БД
