package coingecko

import (
	"context"
	"cryptocurrency/internal/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey, base string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.Wrap(entity.ErrInvalidParams, "CoinGecko client: api key is required")
	}

	u, err := url.Parse(base)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"CoinGecko client: failed to parse URL",
		)
	}
	return &Client{
		baseURL:    u.String(),
		httpClient: &http.Client{Timeout: 5 * time.Second},
		apiKey:     apiKey,
	}, nil
}

func (c *Client) GetRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entity.ErrInvalidParams, "CoinGecko GetRates: titles is required")
	}
	//endpoint, err := c.baseURL.Parse("simple/price")
	rawURL, err := url.Parse(c.baseURL + "/simple/price")
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	query := rawURL.Query()
	query.Set("ids", strings.Join(titles, ","))
	query.Set("vs_currencies", "usd")

	rawURL.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "CoinGecko GetRates: create request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-cg-demo-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "CoinGecko GetRates: request coingecko")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko GetRates: coingecko status code: %s", resp.Status)
	}

	var data map[string]map[string]float64

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "CoinGecko GetRates: decode coingecko response")
	}

	coins := make([]entity.Coin, 0, len(data))
	now := time.Now()

	for title, prices := range data {
		price, ok := prices["usd"]
		if !ok {
			continue
		}
		coin, err := entity.NewCoin(title, price, now)
		if err != nil {
			return nil, errors.Wrap(err, "CoinGecko GetRates: new coin: %w")
		}
		coins = append(coins, *coin)
	}
	if len(coins) == 0 {
		return nil, errors.Wrap(entity.ErrNotFound, "CoinGecko GetRates: coingecko returned no coins")
	}
	return coins, nil
}
