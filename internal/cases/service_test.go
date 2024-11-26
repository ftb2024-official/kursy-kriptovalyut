package cases

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mock_cases "kursy-kriptovalyut/internal/cases/mocks/gen"
	"kursy-kriptovalyut/internal/entities"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name     string
		provider CryptoProvider
		storage  Storage
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
			resErr:   ErrInvalidParam,
		},
		{
			name:     "storage not set",
			provider: mock_cases.NewMockCryptoProvider(ctrl),
			storage:  nil,
			wantErr:  true,
			resErr:   ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service, err := NewService(tt.provider, tt.storage)
			if tt.wantErr {
				require.Nil(t, service)
				require.ErrorIs(t, err, tt.resErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, service)
			require.NotNil(t, service.provider)
			require.NotNil(t, service.storage)
		})
	}
}

func TestHandleMissingTitles(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	mockProvider := mock_cases.NewMockCryptoProvider(ctrl)
	mockStorage := mock_cases.NewMockStorage(ctrl)
	service, _ := NewService(mockProvider, mockStorage)

	tests := []struct {
		name          string
		missingTitles []string
		wantErr       bool
		res           []entities.Coin
		resErr        string
		mockSetup     func()
	}{
		{
			name:          "successful handling",
			missingTitles: []string{"BTC"},
			res: []entities.Coin{
				{
					Title: "BTC",
					Price: 100,
				},
			},
			mockSetup: func() {
				mockProvider.EXPECT().GetActualRates(ctx, []string{"BTC"}).Return([]entities.Coin{
					{
						Title: "BTC",
						Price: 100,
					},
				}, nil)

				mockStorage.EXPECT().Store(ctx, []entities.Coin{
					{
						Title: "BTC",
						Price: 100,
					},
				}).Return(nil)
			},
		},
		{
			name:          "provider failure",
			missingTitles: []string{"BTC"},
			wantErr:       true,
			resErr:        "failed to get coin data from provider",
			mockSetup: func() {
				mockProvider.EXPECT().GetActualRates(ctx, []string{"BTC"}).Return(nil, errors.New("provider error"))
			},
		},
		{
			name:          "storage failure",
			missingTitles: []string{"BTC"},
			wantErr:       true,
			resErr:        "failed to write new coin data to storage",
			mockSetup: func() {
				mockProvider.EXPECT().GetActualRates(ctx, []string{"BTC"}).Return([]entities.Coin{
					{
						Title: "BTC",
						Price: 100,
					},
				}, nil)

				mockStorage.EXPECT().Store(ctx, []entities.Coin{
					{
						Title: "BTC",
						Price: 100,
					},
				}).Return(errors.New("storage error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockSetup()
			res, err := service.handleMissingTitles(ctx, tt.missingTitles)
			if tt.wantErr {
				require.Nil(t, res)
				require.ErrorContains(t, err, tt.resErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.res, res)
		})
	}
}

func TestSplitRequestedTitles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		requested             []string
		existing              []string
		wantExistingReqTitles []string
		wantNewTitles         []string
	}{
		{
			name:                  "no new titles",
			requested:             []string{"BTC", "ETH"},
			existing:              []string{"BTC", "ETH"},
			wantExistingReqTitles: []string{"BTC", "ETH"},
			wantNewTitles:         []string{},
		},
		{
			name:                  "new titles",
			requested:             []string{"USDT", "DOGE"},
			existing:              []string{"BTC", "ETH"},
			wantExistingReqTitles: []string{},
			wantNewTitles:         []string{"USDT", "DOGE"},
		},
		{
			name:                  "no name yet",
			requested:             []string{"BTC", "USDT"},
			existing:              []string{"BTC", "ETH"},
			wantExistingReqTitles: []string{"BTC"},
			wantNewTitles:         []string{"USDT"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotExistingReqTitles, gotNewTitles := splitRequestedTitles(tt.requested, tt.existing)
			require.Equal(t, tt.wantExistingReqTitles, gotExistingReqTitles)
			require.Equal(t, tt.wantNewTitles, gotNewTitles)
		})
	}
}
