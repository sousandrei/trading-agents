package orchestrator

import (
	"context"
	"log/slog"

	"github.com/sousandrei/trading-agents/internal/agents/analysts"
	"github.com/sousandrei/trading-agents/internal/agents/researchers"
	"github.com/sousandrei/trading-agents/internal/agents/risk"
	"github.com/sousandrei/trading-agents/internal/agents/trader"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
	"github.com/sousandrei/trading-agents/internal/types"
)

type Orchestrator struct {
	llm llms.Client
}

func New(llm llms.Client) *Orchestrator {
	return &Orchestrator{
		llm: llm,
	}
}

func (o *Orchestrator) Analyze(
	ctx context.Context,
	position types.Position,
) (*types.Action, error) {
	opts := []llms.GenerateOption{
		// llms.WithDryRun(),
	}

	slog.Info("Running analysis for ticker", "ticker", position.Ticker)
	analystAgents, err := analysts.Run(ctx, o.llm, position, opts...)
	if err != nil {
		return nil, err
	}

	slog.Info("Running research for ticker", "ticker", position.Ticker)
	researchersAgents, err := researchers.Run(ctx, o.llm, analystAgents, position, opts...)
	if err != nil {
		return nil, err
	}

	slog.Info("Running trading strategy for ticker", "ticker", position.Ticker)
	traderAgent, err := trader.Run(ctx, o.llm, analystAgents, researchersAgents, position, opts...)
	if err != nil {
		return nil, err
	}

	slog.Info("Running risk assessment for ticker", "ticker", position.Ticker)
	riskAgents, err := risk.Run(ctx, o.llm, researchersAgents, traderAgent, position, opts...)
	if err != nil {
		return nil, err
	}

	lastMessage := riskAgents["manager"].Messages[len(riskAgents["manager"].Messages)-1].Text
	action, err := types.ParseOutput(position.Ticker, lastMessage)
	if err != nil {
		slog.Error("Failed to parse action", "error", err)
		return nil, err
	}

	slog.Info("Action determined", "action", action.Action, "ticker", position.Ticker)

	return action, nil
}
