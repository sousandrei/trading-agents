package gemini

import (
	"github.com/sousandrei/trading-agents/internal/tools/apiclient"
	"google.golang.org/genai"
)

func (c *Client) apiTools() (map[string]apiclient.Fetch, []*genai.Tool) {
	fns := map[string]apiclient.Fetch{}
	tools := []*genai.Tool{}

	finnhubFns, finnhubTools := c.finnhubTools()
	simfinFns, simfinTools := c.simfinTools()

	tools = append(tools, finnhubTools...)
	tools = append(tools, simfinTools...)

	for _, v := range []map[string]apiclient.Fetch{
		finnhubFns,
		simfinFns,
	} {
		for fnName, fn := range v {
			fns[fnName] = fn
		}
	}

	return fns, tools
}

func searchTool() []*genai.Tool {
	return []*genai.Tool{
		{GoogleSearch: &genai.GoogleSearch{}},
	}
}

func (c *Client) finnhubTools() (map[string]apiclient.Fetch, []*genai.Tool) {
	fns := c.finnhub.GetFunctions()

	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "get_insider_transactions",
					Behavior:    genai.BehaviorBlocking,
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
					Behavior:    genai.BehaviorBlocking,
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
			},
		},
	}

	return fns, tools
}

func (c *Client) simfinTools() (map[string]apiclient.Fetch, []*genai.Tool) {
	fns := c.simfin.GetFunctions()

	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "get_financial_statements",
					Behavior:    genai.BehaviorBlocking,
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
			},
		},
	}

	return fns, tools
}
