package entities

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNewCoin(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		price   float64
		want    *Coin
		wantErr error
	}{
		{"valid input", "ETH", 1000, &Coin{"ETH", 1000, time.Time{}}, nil},
		{"empty title", "", 1000, nil, ErrEmptyTitle},
		{"negative price", "ETH", -1000, nil, ErrNegativeOrZeroPrice},
		{"zero price", "ETH", 0, nil, ErrNegativeOrZeroPrice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coin, err := NewCoin(tt.title, tt.price)
			if !reflect.DeepEqual(coin, tt.want) {
				t.Errorf("got %v, want %v", coin, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAggCoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		title   string
		max     float64
		min     float64
		avg     float64
		want    *AggCoin
		wantErr error
	}{
		{"valid input", "BTC", 1000, 10, 100, &AggCoin{"BTC", 1000, 10, 100}, nil},
		{"empty title", "", 1000, 10, 100, nil, ErrEmptyTitle},
		{"negative price", "ETH", -1000, 10, 100, nil, ErrNegativeOrZeroPrice},
		{"zero price", "ETH", 0, 0, 0, nil, ErrNegativeOrZeroPrice},
		{"min greater than max", "ETH", 100, 1000, 500, nil, ErrMinPriceExceedsMaxOrAvg},
		{"min greater than avg", "ETH", 100, 10, 5, nil, ErrMinPriceExceedsMaxOrAvg},
		{"avg greater than max", "ETH", 100, 10, 200, nil, ErrAvgPriceExceedsMax},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggCoin, err := NewAggCoin(tt.title, tt.max, tt.min, tt.avg)
			if !reflect.DeepEqual(aggCoin, tt.want) {
				t.Errorf("got %v, want %v", aggCoin, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePrices(t *testing.T) {
	tests := []struct {
		name    string
		max     float64
		min     float64
		avg     float64
		wantErr error
	}{
		{"valid input", 1000, 100, 500, nil},
		{"negative price", -1000, -1000, -500, ErrNegativeOrZeroPrice},
		{"zero price", 0, 0, 0, ErrNegativeOrZeroPrice},
		{"min exceeds max or avg", 100, 1000, 500, ErrMinPriceExceedsMaxOrAvg},
		{"avg exceeds max", 100, 100, 500, ErrAvgPriceExceedsMax},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePrices(tt.max, tt.min, tt.avg)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
