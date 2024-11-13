package cases

import (
	"context"
	entity "kursy-kriptovalyut/internal/entities"
)

type Service struct {
	provider CryptoProvider
	storage  Storage
}

func (s *Service) GetLastRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}

func (s *Service) GetMaxRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}

func (s *Service) GetMinRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}

func (s *Service) GetAvgRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	return nil, nil
}

func (s *Service) ActualizeRates(ctx context.Context) error {
	return nil
}

func NewService(provider CryptoProvider, storage Storage) *Service {
	return &Service{provider, storage}
}
