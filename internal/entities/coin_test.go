package entities_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"kursy-kriptovalyut/internal/entities"
)

func TestNewCoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		title   string
		price   float64
		wantErr bool
		resErr  error
	}{
		{
			name:  "valid input",
			title: "ETH",
			price: 1000,
		},
		{
			name:    "empty title",
			title:   "",
			price:   1000,
			wantErr: true,
			resErr:  entities.ErrInvalidParam,
		},
		{
			name:    "wrong price",
			title:   "BTC",
			price:   -1000,
			wantErr: true,
			resErr:  entities.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			coin, err := entities.NewCoin(tt.title, tt.price)
			if tt.wantErr {
				require.Nil(t, coin)
				require.ErrorIs(t, err, tt.resErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.title, coin.Title)
			require.Equal(t, tt.price, coin.Price)
		})
	}
}
