package orchestrator

import (
	"context"

	"github.com/sousandrei/trading-agents/internal/agents/analysts"
	"github.com/sousandrei/trading-agents/internal/agents/researchers"
	"github.com/sousandrei/trading-agents/internal/agents/risk"
	"github.com/sousandrei/trading-agents/internal/agents/trader"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

func Analyze(
	ctx context.Context,
	llm llms.Client,
	ticker string,
) (string, error) {
	opts := []llms.GenerateOption{
		llms.WithDryRun(),
	}

	analystAgents, err := analysts.Run(ctx, llm, ticker, opts...)
	if err != nil {
		return "", err
	}

	researchersAgents, err := researchers.Run(ctx, llm, ticker, analystAgents, opts...)
	if err != nil {
		return "", err
	}

	traderAgent, err := trader.Run(ctx, llm, ticker, analystAgents, researchersAgents, opts...)
	if err != nil {
		return "", err
	}

	riskAgents, err := risk.Run(ctx, llm, ticker, researchersAgents, traderAgent, opts...)
	if err != nil {
		return "", err
	}

	return riskAgents["manager"].Messages[len(riskAgents["manager"].Messages)-1].Text, nil
}
