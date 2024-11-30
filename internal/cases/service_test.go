package cases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"kursy-kriptovalyut/internal/cases"
	mock_cases "kursy-kriptovalyut/internal/cases/mocks/gen"
	"kursy-kriptovalyut/internal/entities"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name     string
		provider cases.CryptoProvider
		storage  cases.Storage
		wantErr  bool
		resErr   error
	}{
		{
			name:     "valid input",
			provider: mock_cases.NewMockCryptoProvider(ctrl),
			storage:  mock_cases.NewMockStorage(ctrl),
		},
		{
			name:     "provider not set",
			provider: nil,
			storage:  mock_cases.NewMockStorage(ctrl),
			wantErr:  true,
			resErr:   entities.ErrInvalidParam,
		},
		{
			name:     "storage not set",
			provider: mock_cases.NewMockCryptoProvider(ctrl),
			storage:  nil,
			wantErr:  true,
			resErr:   entities.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service, err := cases.NewService(tt.provider, tt.storage)
			if tt.wantErr {
				require.Nil(t, service)
				require.ErrorIs(t, err, tt.resErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, service)
		})
	}
}

func TestGetLastRates_Case1_Line_48(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	ctx := context.Background()
	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	requestedTitles := []string{"BTC"}
	btcCoin := entities.Coin{Title: "BTC", Price: 100}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	storage.EXPECT().GetActualCoins(ctx, []string{"BTC"}).Return([]entities.Coin{btcCoin}, nil)

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.NoError(t, err)
	require.ElementsMatch(t, coins, []entities.Coin{btcCoin})
}

func TestGetLastRates_Case2_Line_59(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	ctx := context.Background()
	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	requestedTitles := []string{"ETH"}
	nonExistingCoin := entities.Coin{Title: "ETH", Price: 10}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
	provider.EXPECT().GetActualRates(ctx, requestedTitles).Return([]entities.Coin{nonExistingCoin}, nil)
	storage.EXPECT().Store(ctx, []entities.Coin{nonExistingCoin}).Return(nil)

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.NoError(t, err)
	require.ElementsMatch(t, coins, []entities.Coin{nonExistingCoin})
}

func TestGetLastRates_Case3_Line_73(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	nonExistingCoin := entities.Coin{
		Title: "USDT",
		Price: 10,
	}
	btcCoin := entities.Coin{Title: "BTC", Price: 100}
	ethCoin := entities.Coin{Title: "ETH", Price: 100}
	requestedTitles := []string{"BTC", "ETH", "USDT"}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC", "ETH"}, nil)
	provider.EXPECT().GetActualRates(ctx, []string{"USDT"}).Return([]entities.Coin{nonExistingCoin}, nil)
	storage.EXPECT().Store(ctx, []entities.Coin{nonExistingCoin}).Return(nil)
	storage.EXPECT().GetActualCoins(ctx, []string{"BTC", "ETH"}).Return([]entities.Coin{ethCoin, btcCoin}, nil)

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.NoError(t, err)
	require.ElementsMatch(t, coins, []entities.Coin{nonExistingCoin, btcCoin, ethCoin})
}

func TestGetLastRates_Case4_Line_36(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}

	storage.EXPECT().GetCoinsList(ctx).Return(nil, errors.New("GetCoinsList error"))

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get list of coin titles")
}

func TestGetLastRates_Case5_Line_45(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	storage.EXPECT().GetActualCoins(ctx, []string{"BTC"}).Return(nil, errors.New("GetActualCoins error"))

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get coin data from storage")
}

func TestGetLastRates_Case6_Line_131(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
	provider.EXPECT().GetActualRates(ctx, requestedTitles).Return(nil, errors.New("GetActualRates error"))

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get coin data from provider")
}

func TestGetLastRates_Case7_Line_136(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}
	nonExistingCoin := entities.Coin{Title: "BTC", Price: 100}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
	provider.EXPECT().GetActualRates(ctx, requestedTitles).Return([]entities.Coin{nonExistingCoin}, nil)
	storage.EXPECT().Store(ctx, []entities.Coin{nonExistingCoin}).Return(errors.New("Store error"))

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to write new coin data to storage")
}

func TestGetLastRates_Case8_Line_65(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC", "ETH"}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	provider.EXPECT().GetActualRates(ctx, []string{"ETH"}).Return(nil, errors.New("GetActualRates error"))

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get coin data from provider")
}

func TestGetLastRates_Case9_Line_70(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC", "ETH"}
	nonExistingCoin := entities.Coin{Title: "ETH", Price: 10}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	provider.EXPECT().GetActualRates(ctx, []string{"ETH"}).Return([]entities.Coin{nonExistingCoin}, nil)
	storage.EXPECT().Store(ctx, []entities.Coin{nonExistingCoin}).Return(nil)
	storage.EXPECT().GetActualCoins(ctx, []string{"BTC"}).Return(nil, errors.New("GetActualCoins error"))

	coins, err := srv.GetLastRates(ctx, requestedTitles)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get coin data from storage")
}

func TestGetAggRates_Case1_Line_79(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}
	aggFuncName := "smth"

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "wrong aggregate function name")
}

func TestGetAggRates_Case2_Line_97(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}
	aggFuncName := "max"
	btcCoin := entities.Coin{Title: "BTC", Price: 100}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	storage.EXPECT().GetAggregateCoins(ctx, []string{"BTC"}, aggFuncName).Return([]entities.Coin{btcCoin}, nil)

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.NoError(t, err)
	require.ElementsMatch(t, coins, []entities.Coin{btcCoin})
}

func TestGetAggRates_Case3_Line_108(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}
	aggFuncName := "max"

	storage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
	provider.EXPECT().GetActualRates(ctx, requestedTitles).Return([]entities.Coin{{Title: "BTC", Price: 100}}, nil)
	storage.EXPECT().Store(ctx, []entities.Coin{{Title: "BTC", Price: 100}}).Return(nil)

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.Nil(t, coins)

	expectedErr := fmt.Sprintf("new coins %v added to the storage, but aggregation is unavailable for 5 minutes", requestedTitles)
	require.ErrorContains(t, err, expectedErr)
}

func TestGetAggRates_Case4_Line_124(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC", "ETH"}
	aggFuncName := "max"
	nonExistingCoin := entities.Coin{Title: "ETH", Price: 10}
	existingCoin := entities.Coin{Title: "BTC", Price: 100}

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	provider.EXPECT().GetActualRates(ctx, []string{"ETH"}).Return([]entities.Coin{nonExistingCoin}, nil)
	storage.EXPECT().Store(ctx, []entities.Coin{nonExistingCoin}).Return(nil)
	storage.EXPECT().GetAggregateCoins(ctx, []string{"BTC"}, aggFuncName).Return([]entities.Coin{existingCoin}, nil)

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.ElementsMatch(t, coins, []entities.Coin{existingCoin})

	expectedErr := fmt.Sprintf("partial result returned for coins %v; new coins %v added to the storage, but aggregation is unavailable for 5 minutes", []string{"BTC"}, []string{"ETH"})
	require.ErrorContains(t, err, expectedErr)
}

func TestGetAggRates_Case5_Line_85(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}
	aggFuncName := "max"

	storage.EXPECT().GetCoinsList(ctx).Return(nil, errors.New("ZZZ"))

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get list of coin titles")
}

func TestGetAggRates_Case6_Line_94(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}
	aggFuncName := "max"

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	storage.EXPECT().GetAggregateCoins(ctx, requestedTitles, aggFuncName).Return(nil, errors.New("GetAggregateCoins error"))

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get aggregated coin data from storage")
}

func TestGetAggRates_Case7_Line_105(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC"}
	aggFuncName := "max"

	storage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
	provider.EXPECT().GetActualRates(ctx, requestedTitles).Return(nil, errors.New("GetActualRates"))

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get coin data from provider")
}

func TestGetAggRates_Case8_Line_115(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC", "ETH"}
	aggFuncName := "max"

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	provider.EXPECT().GetActualRates(ctx, []string{"ETH"}).Return(nil, errors.New("GetActualRates error"))

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get coin data from provider")
}

func TestGetAggRates_Case9_Line_120(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := mock_cases.NewMockCryptoProvider(ctrl)
	storage := mock_cases.NewMockStorage(ctrl)

	srv, err := cases.NewService(provider, storage)
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx := context.Background()
	requestedTitles := []string{"BTC", "ETH"}
	aggFuncName := "max"

	storage.EXPECT().GetCoinsList(ctx).Return([]string{"BTC"}, nil)
	provider.EXPECT().GetActualRates(ctx, []string{"ETH"}).Return([]entities.Coin{{Title: "ETH", Price: 10}}, nil)
	storage.EXPECT().Store(ctx, []entities.Coin{{Title: "ETH", Price: 10}}).Return(nil)
	storage.EXPECT().GetAggregateCoins(ctx, []string{"BTC"}, aggFuncName).Return(nil, errors.New("GetAggregateCoins error"))

	coins, err := srv.GetAggRates(ctx, requestedTitles, aggFuncName)
	require.Nil(t, coins)
	require.ErrorContains(t, err, "failed to get aggregated coin data from storage")
}
