package usecases

import (
	"context"
	"kursy-kriptovalyut/internal/entity"
	"kursy-kriptovalyut/internal/ports"
)

type Service struct {
	provider ports.CryptoProvider
	storage  ports.Storage
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

func NewService(provider ports.CryptoProvider, storage ports.Storage) *Service {
	return &Service{provider, storage}
}
