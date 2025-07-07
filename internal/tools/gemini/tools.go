package gemini

import (
	"github.com/sousandrei/trading-agents/internal/tools/apiclient"
	"google.golang.org/genai"
)

func (c *Client) apiTools() (map[string]apiclient.Fetch, []*genai.Tool) {
	finnhubFnMap, finnhubFns := c.finnhubTools()
	simfinFnMap, simfinFns := c.simfinTools()

	fnMap := map[string]apiclient.Fetch{}
	fns := []*genai.FunctionDeclaration{}

	fns = append(fns, finnhubFns...)
	fns = append(fns, simfinFns...)

	for _, v := range []map[string]apiclient.Fetch{
		finnhubFnMap,
		simfinFnMap,
	} {
		for fnName, fn := range v {
			fnMap[fnName] = fn
		}
	}

	return fnMap, []*genai.Tool{
		&genai.Tool{
			FunctionDeclarations: fns,
		},
	}
}

func searchTool() []*genai.Tool {
	return []*genai.Tool{
		{GoogleSearch: &genai.GoogleSearch{}},
	}
}

func (c *Client) finnhubTools() (map[string]apiclient.Fetch, []*genai.FunctionDeclaration) {
	fnMap := c.finnhub.GetFunctions()

	fns := []*genai.FunctionDeclaration{
		{
			Name:        "get_insider_transactions",
			Description: "Fetch insider transactions for a given stock ticker.",
			Parameters: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"ticker": {
						Type:        "string",
						Description: "The ticker symbol of the stock to query.",
					},
				},
			},
		},
		{
			Name:        "get_insider_sentiment",
			Description: "Fetch insider sentiment for a given stock ticker.",
			Parameters: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"ticker": {
						Type:        "string",
						Description: "The ticker symbol of the stock to query.",
					},
				},
			},
		},
	}

	return fnMap, fns
}

func (c *Client) simfinTools() (map[string]apiclient.Fetch, []*genai.FunctionDeclaration) {
	fnMap := c.simfin.GetFunctions()

	fns := []*genai.FunctionDeclaration{
		{
			Name:        "get_financial_statements",
			Description: "Fetch financial statements for a given stock ticker.",
			Parameters: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"ticker": {
						Type:        "string",
						Description: "The ticker symbol of the stock to query.",
					},
				},
			},
		},
	}

	return fnMap, fns
}
