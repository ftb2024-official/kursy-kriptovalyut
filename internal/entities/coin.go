package entities

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrSmthWentWrong = errors.New("smth went wrong")
	ErrEmptyTitle    = errors.New("empty title")
	ErrNegativePrice = errors.New("negative price")
	ErrZeroPrice     = errors.New("zero price")
)

type Coin struct {
	title    string
	price    float64
	actualAt time.Time
}

func NewCoin(title string, price float64) (*Coin, error) {
	if title == "" {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrEmptyTitle)
	}

	if price < 0 {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrNegativePrice)
	}

	if price == 0 {
		return nil, fmt.Errorf("%w: %w", ErrSmthWentWrong, ErrZeroPrice)
	}

	return &Coin{
		title: title,
		price: price,
	}, nil
}

func (c *Coin) SetTime() {
	c.actualAt = time.Now()
}
