package orchestrator

import (
	"context"
	"log/slog"

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
		// llms.WithDryRun(),
	}

	slog.Info("Running analysis for ticker", "ticker", ticker)
	analystAgents, err := analysts.Run(ctx, llm, ticker, opts...)
	if err != nil {
		return "", err
	}

	slog.Info("Running research for ticker", "ticker", ticker)
	researchersAgents, err := researchers.Run(ctx, llm, ticker, analystAgents, opts...)
	if err != nil {
		return "", err
	}

	slog.Info("Running trading strategy for ticker", "ticker", ticker)
	traderAgent, err := trader.Run(ctx, llm, ticker, analystAgents, researchersAgents, opts...)
	if err != nil {
		return "", err
	}

	slog.Info("Running risk assessment for ticker", "ticker", ticker)
	riskAgents, err := risk.Run(ctx, llm, ticker, researchersAgents, traderAgent, opts...)
	if err != nil {
		return "", err
	}

	return riskAgents["manager"].Messages[len(riskAgents["manager"].Messages)-1].Text, nil
}
