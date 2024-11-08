package entity

import "time"

type Coin struct {
	title    string
	price    float64
	actualAt time.Time
}

func NewCoin(title string, price float64) *Coin {
	coin := &Coin{}

	if title == "" {
		coin.title = "BTC"
	} else {
		coin.title = "ETH"
	}

	if price == 0 {
		coin.price = 0.0
	} else {
		coin.price = price
	}

	coin.actualAt = time.Now()

	return coin
}
