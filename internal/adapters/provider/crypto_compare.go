package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"kursy-kriptovalyut/internal/entities"
	"kursy-kriptovalyut/pkg/logger"
)

var log = logger.NewLogger()

type CryptoCompare struct {
	baseUrl    string
	apiKey     string
	httpClient *http.Client
}

func NewCryptoCompare(baseUrl string, apiKey string) (*CryptoCompare, error) {
	if baseUrl == "" || apiKey == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "base-url/api-key is empty")
	}

	return &CryptoCompare{
		baseUrl:    baseUrl,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

const (
	fromSyms = "fsyms"
	toSyms   = "tsyms"
	currency = "USD"
	max      = "MAX"
	min      = "MIN"
	avg      = "AVG"
)

func (cc *CryptoCompare) GetActualRates(ctx context.Context, titles []string, extraArg string) ([]entities.Coin, error) {
	rawURL, err := url.Parse(cc.baseUrl)
	if err != nil {
		log.Error("(GetActualRates) failed to parse base url:", zap.Any("url", cc.baseUrl))
		return nil, errors.Wrapf(entities.ErrInternal, "failed to parse url: %v", err)
	}

	// manual raw query
	// rawQuery := fmt.Sprintf("%s=%s&%s=%s", fromSyms, strings.Join(titles, ","), toSyms, currency)
	// rawURL.RawQuery = rawQuery

	queries := rawURL.Query()
	queries.Add(fromSyms, strings.Join(titles, ","))
	queries.Add(toSyms, currency)

	rawURL.RawQuery, err = url.QueryUnescape(queries.Encode())
	if err != nil {
		log.Error("(GetActualRates) failed to encode url:", zap.Any("url", rawURL.String()))
		return nil, errors.Wrapf(entities.ErrInternal, "failed to encode url: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL.String(), nil)
	if err != nil {
		log.Error("(GetActualRates) failed to create new request")
		return nil, errors.Wrapf(entities.ErrInternal, "failed to create new request, err: %v", err)
	}

	req.Header.Set("Authorization", "Apikey "+cc.apiKey)

	resp, err := cc.httpClient.Do(req)
	if err != nil {
		log.Error("(GetActualRates) failed to execute request")
		return nil, errors.Wrapf(entities.ErrInternal, "failed to execute request, err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Info("(GetActualRates) unexpected status code:", zap.Any("statusCode", resp.StatusCode))
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("(GetActualRates) failed to read response body")
		return nil, errors.Wrap(err, "failed to read response body")
	}

	strBody := string(body)
	if strings.Contains(strBody, `"Response":"Error"`) {
		log.Error("(GetActualRates) non-existing coin(s):", zap.Any("coinTitle(s)", titles))
		return nil, errors.Wrapf(entities.ErrNotFound, "coin %v does not exist", titles)
	}

	type CryptoData struct {
		RAW map[string]struct {
			USD struct {
				PRICE   float64 `json:"PRICE"`
				HIGHDAY float64 `json:"HIGHDAY"`
				LOWDAY  float64 `json:"LOWDAY"`
			}
		}
	}

	var data CryptoData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("(GetActualRates) failed to parse resp.body")
		return nil, errors.Wrap(err, "failed to parse response body, invalid JSON format")
	}

	coins := make([]entities.Coin, 0, len(data.RAW))
	var price float64
	for coinTitle, info := range data.RAW {
		switch extraArg {
		case max:
			price = info.USD.HIGHDAY
		case min:
			price = info.USD.LOWDAY
		case avg:
			price = (info.USD.PRICE + info.USD.HIGHDAY + info.USD.LOWDAY) / 3
		default:
			price = info.USD.PRICE
		}

		coin, err := entities.NewCoin(coinTitle, price)
		if err != nil {
			log.Error("(GetActualRates) failed to create new coin:", zap.Any("newCoin", coin))
			return nil, errors.Wrap(err, "failed to create new coin")
		}

		coins = append(coins, *coin)
	}

	return coins, nil
}
