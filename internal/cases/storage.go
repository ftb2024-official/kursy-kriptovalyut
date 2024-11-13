package cases

import (
	"context"
	entity "kursy-kriptovalyut/internal/entities"
)

type Storage interface {
	Store(ctx context.Context, coins []entity.Coin) error
	GetCoinsList(ctx context.Context) ([]string, error)
	GetActualCoins(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetAggregateCoins(ctx context.Context, titles []string) ([]entity.Coin, error)
}

// Store - для записи новых данных в БД полученных из внешнего API
// GetCoinsList - для получения списка коинов
// GetActualCoins - для получения последней цены крипты записанной в БД
// GetAggregateCoins - ???
