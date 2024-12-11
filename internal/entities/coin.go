package entities

import (
	"github.com/pkg/errors"
)

type Coin struct {
	Title string
	Price float64
}

func NewCoin(title string, price float64) (*Coin, error) {
	if title == "" {
		return nil, errors.Wrap(ErrInvalidParam, "title must not be empty")
	}

	if price <= 0 {
		return nil, errors.Wrap(ErrInvalidParam, "price must be positive")
	}

	return &Coin{
		Title: title,
		Price: price,
	}, nil
}
