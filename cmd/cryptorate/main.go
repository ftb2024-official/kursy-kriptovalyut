package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// type Coin struct {
// 	Title string
// 	Price float64
// }

// const (
// 	fromSyms = "fsyms"
// 	toSyms   = "tsyms"
// 	currency = "USD"
// )

// func GetActualRates(ctx context.Context, baseUrl, apiKey string, titles []string) ([]Coin, error) {
// 	URLRaw, err := url.Parse(baseUrl)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse url: %v", err)
// 	}

// 	// manual raw query
// 	// rawQuery := fmt.Sprintf("%s=%s&%s=%s", fromSyms, strings.Join(titles, ","), toSyms, currency)
// 	// fUrl.RawQuery = rawQuery

// 	queries := URLRaw.Query()
// 	queries.Add(fromSyms, strings.Join(titles, ","))
// 	queries.Add(toSyms, currency)

// 	URLRaw.RawQuery = queries.Encode()
// 	if err != nil {
// 		return nil, errors.Wrapf(entities.ErrInternal, "failed to encode url: %v", err)
// 	}

// 	fmt.Println(URLRaw.String())

// req, err := http.NewRequestWithContext(ctx, http.MethodGet, fUrl.String(), nil)
// if err != nil {
// 	return nil, fmt.Errorf("failed to create new request, err: %v", err)
// }

// req.Header.Set("Authorization", "Apikey "+apiKey)

// client := http.Client{}

// resp, err := client.Do(req)
// if err != nil {
// 	return nil, fmt.Errorf("failed to execute request, err: %v", err)
// }
// defer resp.Body.Close()

// if resp.StatusCode != http.StatusOK {
// 	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// }

// body, err := io.ReadAll(resp.Body)
// if err != nil {
// 	return nil, fmt.Errorf("failed to read response body: %v", err)
// }

// var data map[string]map[string]float64
// err = json.Unmarshal(body, &data)
// if err != nil {
// 	return nil, fmt.Errorf("failed to parse response body, invalid JSON format: %v", err)
// }

// coins := make([]Coin, 0, len(data))
// for title, value := range data {
// 	price := value[currency]

// 	coin := Coin{Title: title, Price: price}

// 	coins = append(coins, coin)
// }
// return nil, nil
// }

// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()

// 	titles := []string{"BTC", "ETH", "DOGE"}
// 	url := "https://min-api.cryptocompare.com/data/pricemulti"
// 	apiKey := "851e396ad68e892830b474f074b051d2104b77576c25b9058ef16d4a477515d8"

// 	_, _ = GetActualRates(ctx, url, apiKey, titles)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// current := time.Now().Format("2006-01-02")

// 	// fmt.Println(coins)
// 	// fmt.Println(current)
// }

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/jackc/pgx/v5"
// )

func main() {
	connStr := "postgresql://user:pswd@localhost:5434/crypto_rate?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected...")

}
