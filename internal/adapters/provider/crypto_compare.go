package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"kursy-kriptovalyut/internal/entities"
)

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
)

func (cc *CryptoCompare) GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	rawURL, err := url.Parse(cc.baseUrl)
	if err != nil {
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
		return nil, errors.Wrapf(entities.ErrInternal, "failed to encode url: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL.String(), nil)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to create new request, err: %v", err)
	}

	req.Header.Set("Authorization", "Apikey "+cc.apiKey)

	resp, err := cc.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to execute request, err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var data map[string]map[string]float64
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse response body, invalid JSON format")
	}

	coins := make([]entities.Coin, 0, len(data))
	for title, value := range data {
		price := value[currency]

		coin, err := entities.NewCoin(title, price)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create new coin")
		}

		coins = append(coins, *coin)
	}
	return coins, nil
}
