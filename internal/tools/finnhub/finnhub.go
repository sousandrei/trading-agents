package finnhub

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sousandrei/trading-agents/internal/tools/apiclient"
)

const baseURL = "https://finnhub.io/api/v1"

type Client struct {
	apiClient *apiclient.Client
	apiKey    string
}

func New(
	apiKey string,
	opts ...apiclient.Options,
) (*Client, error) {
	apiClient, err := apiclient.New(opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		apiClient: apiClient,
		apiKey:    apiKey,
	}, nil
}

func (c *Client) GetFunctions() map[string]apiclient.Fetch {
	return map[string]apiclient.Fetch{
		"get_company_profile":      c.getCompanyProfile,
		"get_insider_transactions": c.getInsiderTransactions,
		"get_insider_sentiment":    c.getInsiderSentiment,
	}
}

func (c *Client) getCompanyProfile(
	ctx context.Context,
	args map[string]any,
) (map[string]any, error) {
	symbol, ok := args["ticker"].(string)
	if !ok || symbol == "" {
		return nil, fmt.Errorf("ticker symbol is required")
	}

	reqURL, _ := url.Parse(baseURL + "/stock/profile2")

	query := reqURL.Query()
	query.Set("symbol", symbol)

	reqURL.RawQuery = query.Encode()

	headers := map[string]string{
		"X-Finnhub-Token": c.apiKey,
	}

	var res map[string]any
	err := c.apiClient.Get(ctx, reqURL.String(), headers, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) getInsiderTransactions(
	ctx context.Context,
	args map[string]any,
) (map[string]any, error) {
	symbol, ok := args["ticker"].(string)
	if !ok || symbol == "" {
		return nil, fmt.Errorf("ticker symbol is required")
	}

	reqURL, _ := url.Parse(baseURL + "/stock/insider-transactions")

	query := reqURL.Query()
	query.Set("symbol", symbol)

	reqURL.RawQuery = query.Encode()

	headers := map[string]string{
		"X-Finnhub-Token": c.apiKey,
	}

	var res map[string]any
	err := c.apiClient.Get(ctx, reqURL.String(), headers, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) getInsiderSentiment(
	ctx context.Context,
	args map[string]any,
) (map[string]any, error) {
	symbol, ok := args["ticker"].(string)
	if !ok || symbol == "" {
		return nil, fmt.Errorf("ticker symbol is required")
	}

	reqURL, _ := url.Parse(baseURL + "/stock/insider-sentiment")

	query := reqURL.Query()
	query.Set("symbol", symbol)

	reqURL.RawQuery = query.Encode()

	headers := map[string]string{
		"X-Finnhub-Token": c.apiKey,
	}

	var res map[string]any
	err := c.apiClient.Get(ctx, reqURL.String(), headers, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
