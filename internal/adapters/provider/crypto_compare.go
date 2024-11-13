package provider

import (
	"context"
	entity "kursy-kriptovalyut/internal/entities"
)

type CryptoCompare struct {
	apiKey string
}

func NewCryptoCompare(apiKey string) *CryptoCompare {
	return &CryptoCompare{apiKey: apiKey}
}

func (cc *CryptoCompare) GetActualRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}
