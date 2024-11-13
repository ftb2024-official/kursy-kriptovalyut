package entities

import "time"

type Coin struct {
	title    string
	price    float64
	actualAt time.Time
}

func NewCoin(title string, price float64) (*Coin, error) {
	if title == "" {
		return nil, nil
	}

	if price <= 0 {
		price = 0
	}

	return &Coin{
		title: title,
		price: price,
	}, nil
}

func (c *Coin) SetTime() {
	c.actualAt = time.Now()
}
