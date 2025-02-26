package cases

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"kursy-kriptovalyut/internal/entities"
	"kursy-kriptovalyut/pkg/logger"
)

var log = logger.NewLogger()

type Service struct {
	provider CryptoProvider
	storage  Storage
}

func NewService(provider CryptoProvider, storage Storage) (*Service, error) {
	if provider == nil || provider == CryptoProvider(nil) {
		return nil, errors.Wrap(entities.ErrInvalidParam, "provider not set")
	}

	if storage == nil || storage == Storage(nil) {
		return nil, errors.Wrap(entities.ErrInvalidParam, "storage not set")
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
		log.Error("(service.GetLastRates) failed to get list of coin titles")
		return nil, errors.Wrap(err, "failed to get list of coin titles")
	}

	existingReqTitles, nonExistingReqTitles := splitRequestedTitles(requestedCoinTitles, existingTitles)

	// 1-случай: все монеты есть в хранилище
	if len(nonExistingReqTitles) == 0 {
		log.Info("(service.GetLastRates) all requested titles in storage")
		coins, err := s.storage.GetActualCoins(ctx, requestedCoinTitles)
		if err != nil {
			log.Error("(service.GetLastRates) failed to get coin data from storage")
			return nil, errors.Wrap(err, "failed to get coin data from storage")
		}

		return coins, nil
	}

	// 2-случай: все запрашиваемые монеты отсутствуют в хранилище
	if len(existingReqTitles) == 0 {
		log.Info("(service.GetLastRates) all requested titles not in storage")
		newCoins, err := s.handleMissingTitles(ctx, nonExistingReqTitles, "PRICE")
		if err != nil {
			log.Error("(service.GetLastRates) error from (handleMissingTitles)")
			return nil, err
		}

		return newCoins, nil
	}

	// 3-случай: часть монет есть в хранилище, часть отсутствует
	log.Info("(service.GetLastRates) requested titles partially in storage")
	newCoins, err := s.handleMissingTitles(ctx, nonExistingReqTitles, "PRICE")
	if err != nil {
		log.Error("(service.GetLastRates) error from (handleMissingTitles)")
		return nil, err
	}

	coins, err := s.storage.GetActualCoins(ctx, existingReqTitles)
	if err != nil {
		log.Error("(service.GetLastRates) failed to get coin data from storage")
		return nil, errors.Wrap(err, "failed to get coin data from storage")
	}

	return append(coins, newCoins...), nil
}

func (s *Service) GetAggRates(ctx context.Context, requestedCoinTitles []string, aggFuncName string) ([]entities.Coin, error) {
	validAggFuncs := map[string]bool{"MAX": true, "MIN": true, "AVG": true}
	if !validAggFuncs[strings.ToUpper(aggFuncName)] {
		log.Error("(service.GetAggRates) wrong aggregate function")
		return nil, errors.Wrap(entities.ErrInvalidParam, "wrong aggregate function name")
	}

	// получаем список монет, которые уже есть в хранилище
	existingTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		log.Error("(service.GetAggRates) failed to get list of coin titles")
		return nil, errors.Wrap(err, "failed to get list of coin titles")
	}

	existingReqTitles, nonExistingReqTitles := splitRequestedTitles(requestedCoinTitles, existingTitles)

	// 1-случай: все запрашиваемые монеты есть в хранилище
	if len(nonExistingReqTitles) == 0 {
		log.Info("(service.GetAggRates) all requested titles in storage")
		aggCoins, err := s.storage.GetAggregateCoins(ctx, requestedCoinTitles, aggFuncName)
		if err != nil {
			log.Error("(service.GetAggRates) failed to get aggregated coin data from storage")
			return nil, errors.Wrap(err, "failed to get aggregated coin data from storage")
		}

		return aggCoins, nil
	}

	// 2-случай: все запрашиваемые монеты отсутствуют в хранилище
	if len(existingReqTitles) == 0 {
		log.Info("(service.GetAggRates) all requested titles not in storage")
		// получаем актуальные данные по отсутствующим монетам от провайдера и сохраняем в хранилище
		newAggCoins, err := s.handleMissingTitles(ctx, nonExistingReqTitles, aggFuncName)
		if err != nil {
			log.Error("(service.GetAggRates) error from (handleMissingTitles)")
			return nil, err
		}

		return newAggCoins, nil
	}

	// 3-случай: часть монет есть в хранилище, часть отсутствует
	// получаем актуальные данные по отсутствующим монетам от провайдера
	log.Info("(service.GetAggRates) requested titles partially in storage")
	newAggCoins, err := s.handleMissingTitles(ctx, nonExistingReqTitles, aggFuncName)
	if err != nil {
		log.Error("(service.GetAggRates) error from (handleMissingTitles)")
		return nil, err
	}

	aggCoins, err := s.storage.GetAggregateCoins(ctx, existingReqTitles, aggFuncName)
	if err != nil {
		log.Error("(service.GetAggRates) failed to get aggregated coin data from storage")
		return nil, errors.Wrap(err, "failed to get aggregated coin data from storage")
	}

	// возвращаем частичный результат и ошибку
	return append(aggCoins, newAggCoins...), nil
}

func (s *Service) ActualizeRates(ctx context.Context) error {
	log.Info("(service.ActualizeRates) getting list of coin titles")
	existingTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		log.Error("(service.ActualizeRates) failed to get list of coin titles")
		return errors.Wrap(err, "failed to get list of coin titles")
	}

	if len(existingTitles) == 0 {
		log.Info("(service.ActualizeRates) no coin titles in storage yet")
		return nil
	}

	log.Info("(service.ActualizeRates) actualizing coin rates")
	_, err = s.handleMissingTitles(ctx, existingTitles, "PRICE")
	if err != nil {
		log.Error("(service.ActualizeRates) failed to actualize coin rates")
		return errors.Wrap(err, "failed to actualize coin rates")
	}

	return nil
}

func (s *Service) handleMissingTitles(ctx context.Context, missingTitles []string, extraArg string) ([]entities.Coin, error) {
	// получаем актуальные данные по отсутствующим монетам от провайдера
	newCoins, err := s.provider.GetActualRates(ctx, missingTitles, extraArg)
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
