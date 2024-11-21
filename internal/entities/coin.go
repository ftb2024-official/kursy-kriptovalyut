package entities

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrSmthWentWrong           = errors.New("smth went wrong")
	ErrEmptyTitle              = errors.New("empty title")
	ErrNegativeOrZeroPrice     = errors.New("negative or zero price")
	ErrMinPriceExceedsMaxOrAvg = errors.New("min price greater than max or avg price")
	ErrAvgPriceExceedsMax      = errors.New("avg price greater than max price")
)

type Coin struct {
	title    string
	price    float64
	actualAt time.Time
}

type AggCoin struct {
	title    string
	maxPrice float64
	minPrice float64
	avgPrice float64
}

func NewCoin(title string, price float64) (*Coin, error) {
	if title == "" {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrEmptyTitle)
	}

	if price <= 0 {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrNegativeOrZeroPrice)
	}

	return &Coin{
		title: title,
		price: price,
	}, nil
}

func NewAggCoin(title string, max, min, avg float64) (*AggCoin, error) {
	if title == "" {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrEmptyTitle)
	}

	if err := validatePrices(max, min, avg); err != nil {
		return nil, err
	}

	return &AggCoin{
		title:    title,
		maxPrice: max,
		minPrice: min,
		avgPrice: avg,
	}, nil
}

func validatePrices(max, min, avg float64) error {
	if max <= 0 || min <= 0 || avg <= 0 {
		return fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrNegativeOrZeroPrice)
	}

	if min > max || min > avg {
		return fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrMinPriceExceedsMaxOrAvg)
	}

	if avg > max {
		return fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrAvgPriceExceedsMax)
	}

	return nil
}

func (c *Coin) SetTime() {
	c.actualAt = time.Now()
}
