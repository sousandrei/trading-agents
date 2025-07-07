package simfin

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/sousandrei/trading-agents/internal/tools/apiclient"
)

const baseURL = "https://backend.simfin.com/api/v3"

type Client struct {
	apiClient *apiclient.Client
	apiKey    string
}

func New(
	apiKey string,
) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) GetFunctions() map[string]apiclient.Fetch {
	return map[string]apiclient.Fetch{
		"get_financial_statements": c.getFinancialStatements,
	}

}

// https://backend.simfin.com/api/v3/companies/statements/compact?
// ticker=NVDA
// &statements=PL,BS,CF,DERIVED
// &period=Q1,Q2,Q3,Q4
// &start=2024-06-01
// &end=2025-06-01'

func (c *Client) getFinancialStatements(
	ctx context.Context,
	args map[string]any,
) (map[string]any, error) {
	symbol, ok := args["ticker"].(string)
	if !ok || symbol == "" {
		return nil, fmt.Errorf("ticker symbol is required")
	}

	var endDate string
	if end, ok := args["endDate"].(string); ok && end != "" {
		endDate = end
	} else {
		endDate = time.Now().Format("2006-01-02")
	}

	var startDate string
	if start, ok := args["startDate"].(string); ok && start != "" {
		startDate = start
	} else {
		startDate = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	}

	reqURL, _ := url.Parse(baseURL + "/companies/statements/compact")

	query := reqURL.Query()
	query.Set("ticker", symbol)
	query.Set("statements", "PL,BS,CF,DERIVED")
	query.Set("period", "Q1,Q2,Q3,Q4")
	query.Set("start", startDate)
	query.Set("end", endDate)

	headers := map[string]string{
		"Authorization": "api-key " + c.apiKey,
	}

	var res map[string]any
	err := c.apiClient.Get(ctx, reqURL.String(), headers, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
