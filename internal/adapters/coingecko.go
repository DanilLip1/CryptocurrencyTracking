package adapters

import (
	"context"
	"cryptocurrency/internal/entity"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string, httpClient *http.Client) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("api key is required")
	}
	u, err := url.Parse("https://api.coingecko.com/api/v3/")
	if err != nil {
		return nil, err
	}
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 5,
		}
	}
	return &Client{
		baseURL:    u,
		httpClient: httpClient,
		apiKey:     apiKey,
	}, nil
}
func (c *Client) GetRates(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("titles is empty")
	}
	endpoint, err := c.baseURL.Parse("simple/price")
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	query := url.Values{}
	query.Set("ids", strings.Join(titles, ","))
	query.Set("vs_currencies", "usd")

	endpoint.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-cg-demo-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request coingecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coingecko status code: %s", resp.Status)
	}

	var data map[string]map[string]float64

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode coingecko response: %w", err)
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
			return nil, fmt.Errorf("new coin: %w", err)
		}
		coins = append(coins, *coin)
	}
	if len(coins) == 0 {
		return nil, fmt.Errorf("coingecko returned no coins")
	}
	return coins, nil
}
