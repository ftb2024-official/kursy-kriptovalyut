package dto

type CoinDTO struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type ErrRespDTO map[string]string
