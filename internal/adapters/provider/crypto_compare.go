package provider

import (
	"context"
	"kursy-kriptovalyut/internal/entity"
)

type CryptoCompare struct{}

func (cc *CryptoCompare) GetActualRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}
