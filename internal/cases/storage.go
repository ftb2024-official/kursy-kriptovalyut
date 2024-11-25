package cases

import (
	"context"

	"kursy-kriptovalyut/internal/entities"
)

//go:generate mockgen -source=./storage.go -destination=./mocks/gen/mock_storage.go
type Storage interface {
	Store(ctx context.Context, coins []entities.Coin) error
	GetCoinsList(ctx context.Context) ([]string, error)
	GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAggregateCoins(ctx context.Context, titles []string, aggFuncName string) ([]entities.Coin, error)
}

// Store - для записи новых данных в БД полученных из внешнего API
// GetCoinsList - для получения списка коинов
// GetActualCoins - для получения последней цены крипты записанной в БД
// GetAggregateCoins - для получения агрегированного ответа (макс. мин. сред. цены)
