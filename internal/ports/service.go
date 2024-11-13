package ports

import (
	"context"
	entity "kursy-kriptovalyut/internal/entities"
)

type Service interface {
	GetLastRates(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMaxRates(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMinRates(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetAvgRates(ctx context.Context, titles []string) ([]entity.Coin, error)
}

// GetLastRates - для получения актуального курса
// GetMaxRates - для получения макс.цены крипты за день
// GetMinRates - для получения мин.цены крипты за день
// GetAvgRates - для получения ср.цены крипты за день
