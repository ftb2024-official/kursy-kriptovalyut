package cases

import (
	"context"
	"errors"
	"fmt"
	entity "kursy-kriptovalyut/internal/entities"
)

var (
	ErrSmthWentWrong = errors.New("smth went wrong")
	ErrNilProvider   = errors.New("nil provider")
	ErrNilStorage    = errors.New("nil storage")
)

type Service struct {
	provider CryptoProvider
	storage  Storage
}

func NewService(provider CryptoProvider, storage Storage) (*Service, error) {
	if provider == nil || provider == CryptoProvider(nil) {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrNilProvider)
	}

	if storage == nil || storage == Storage(nil) {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrNilStorage)
	}

	return &Service{
		provider: provider,
		storage:  storage,
	}, nil
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
