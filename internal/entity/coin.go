package entity

import "time"

type Coin struct {
	title    string
	price    float64
	actualAt time.Time
}

func NewCoin(title string, price float64) *Coin {
	if title != "BTC" && title != "ETH" {
		title = "BTC"
	}

	if price < 0 {
		price = 0
	}

	return &Coin{title, price, time.Now()}
}
