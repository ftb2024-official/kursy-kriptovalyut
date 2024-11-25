package cases

import (
	"context"

	"github.com/pkg/errors"

	"kursy-kriptovalyut/internal/entities"
)

type Service struct {
	provider CryptoProvider
	storage  Storage
}

func NewService(provider CryptoProvider, storage Storage) (*Service, error) {
	if provider == nil || provider == CryptoProvider(nil) {
		return nil, errors.Wrap(ErrInvalidParam, "provider not set")
	}

	if storage == nil || storage == Storage(nil) {
		return nil, errors.Wrap(ErrInvalidParam, "storage not set")
	}

	return &Service{
		provider: provider,
		storage:  storage,
	}, nil
}

func (s *Service) GetLastRates(ctx context.Context, requestedCoinTitles []string) ([]entities.Coin, error) {
	// получаем список монет, которые уже есть в хранилище
	existingTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of coin titles")
	}

	existingReqTitles, nonExistingReqTitles := splitRequestedTitles(requestedCoinTitles, existingTitles)

	// 1-случай: все монеты есть в хранилище
	if len(nonExistingReqTitles) == 0 {
		coins, err := s.storage.GetActualCoins(ctx, requestedCoinTitles)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get coin data from storage")
		}

		return coins, nil
	}

	// 2-случай: все запрашиваемые монеты отсутствуют в хранилище
	if len(existingReqTitles) == 0 {
		newCoins, err := s.handleMissingTitles(ctx, nonExistingReqTitles)
		if err != nil {
			// если ошибка от провайдера (nil, err), если от хранилища (newCoins, nil)
			return nil, err
		}

		return newCoins, nil
	}

	// 3-случай: часть монет есть в хранилище, часть отсутствует
	newCoins, err := s.handleMissingTitles(ctx, nonExistingReqTitles)
	if err != nil {
		return nil, err
	}

	coins, err := s.storage.GetActualCoins(ctx, existingReqTitles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get coin data from storage")
	}

	return append(coins, newCoins...), nil
}

func (s *Service) GetAggRates(ctx context.Context, requestedCoinTitles []string, aggFuncName string) ([]entities.Coin, error) {
	// получаем список монет, которые уже есть в хранилище
	existingTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of coin titles")
	}

	existingReqTitles, nonExistingReqTitles := splitRequestedTitles(requestedCoinTitles, existingTitles)

	// 1-случай: все запрашиваемые монеты есть в хранилище
	if len(nonExistingReqTitles) == 0 {
		aggCoins, err := s.storage.GetAggregateCoins(ctx, requestedCoinTitles, aggFuncName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get aggregated coin data from storage")
		}

		return aggCoins, nil
	}

	// 2-случай: все запрашиваемые монеты отсутствуют в хранилище
	if len(existingReqTitles) == 0 {
		// получаем актуальные данные по отсутствующим монетам от провайдера
		_, err := s.handleMissingTitles(ctx, nonExistingReqTitles)
		if err != nil {
			return nil, err
		}

		return nil, errors.Wrapf(err, "new coins %v added to the storage, but aggregation is unavailable for 5 minutes", nonExistingReqTitles)
	}

	// 3-случай: часть монет есть в хранилище, часть отсутствует
	aggCoins, err := s.storage.GetAggregateCoins(ctx, existingReqTitles, aggFuncName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aggregated coin data from storage")
	}

	// получаем актуальные данные по отсутствующим монетам от провайдера
	_, err = s.handleMissingTitles(ctx, nonExistingReqTitles)
	if err != nil {
		return nil, err
	}

	// возвращаем частичный результат и ошибку
	return aggCoins, errors.Wrapf(err, "partial result returned for coins %v; new coins %v added to the storage, but aggregation is unavailable for 5 minutes", existingReqTitles, nonExistingReqTitles)
}

func (s *Service) handleMissingTitles(ctx context.Context, missingTitles []string) ([]entities.Coin, error) {
	// получаем актуальные данные по отсутствующим монетам от провайдера
	newCoins, err := s.provider.GetActualRates(ctx, missingTitles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get coin data from provider")
	}

	// сохраняем новые монеты в хранилище
	if err := s.storage.Store(ctx, newCoins); err != nil {
		return nil, errors.Wrap(err, "failed to write new coin data to storage")
	}

	return newCoins, nil
}

// функция для разделения монет на категории: (существующий запрашиваемый) и (несуществующий запрашиваемый)
func splitRequestedTitles(requested, existing []string) ([]string, []string) {
	existingReqTitles := make([]string, 0)
	nonExistingReqTitles := make([]string, 0)

	existingTitlesMap := make(map[string]struct{}, len(existing))
	for _, title := range existing {
		existingTitlesMap[title] = struct{}{}
	}

	for _, title := range requested {
		if _, ok := existingTitlesMap[title]; ok {
			existingReqTitles = append(existingReqTitles, title)
		} else {
			nonExistingReqTitles = append(nonExistingReqTitles, title)
		}
	}
	return existingReqTitles, nonExistingReqTitles
}
