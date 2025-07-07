package trader

import (
	"context"
	"fmt"

	"github.com/sousandrei/trading-agents/internal/agents"
	"github.com/sousandrei/trading-agents/internal/agents/analysts"
	"github.com/sousandrei/trading-agents/internal/agents/researchers"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

const traderPrompt = `You are a trading agent analyzing market data to make investment decisions.
Based on your analysis, provide a specific recommendation to buy, sell, or hold.
End with a firm decision and always conclude your response with 'FINAL TRANSACTION PROPOSAL: **BUY/HOLD/SELL**' to confirm your recommendation.`

func Run(
	ctx context.Context,
	llm llms.Client,
	ticker string,
	analystAgents map[string]agents.Agent,
	researcherAgents map[string]agents.Agent,
	opts ...llms.GenerateOption,
) (*agents.Agent, error) {
	if opts == nil {
		opts = []llms.GenerateOption{}
	}

	prompt := fmt.Sprintf("%s\nStock in question: %s", traderPrompt, ticker)
	prompt = analysts.AppendOutput(prompt, analystAgents)
	prompt = researchers.AppendManagerOutput(prompt, researcherAgents)

	res, err := llm.Generate(ctx, prompt, opts...)
	if err != nil {
		return nil, fmt.Errorf("error running manager: %w", err)
	}

	// TODO: remove, dev only
	agents.WriteMessagesToFile("trader", "trader", res)

	return &agents.Agent{
		Prompt:   prompt,
		Messages: res,
	}, nil
}

func AppendOutput(prompt string, trader *agents.Agent) string {
	prompt += fmt.Sprintf("%s\n\n\n\nTrader report:", prompt)
	return prompt
}
