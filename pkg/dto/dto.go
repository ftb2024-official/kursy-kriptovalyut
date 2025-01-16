package dto

type CoinDTO struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type ErrRespDTO struct {
	StatusCode int    `json:"status_code"`
	Msg        string `json:"msg"`
}

// to struct ✅
// error code 400 ✅
// config file ✅
// cron job
// logger peresmotret
// application layer ✅
// Dockerfile
