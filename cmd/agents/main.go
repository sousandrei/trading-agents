package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/caarlos0/env"

	"github.com/sousandrei/trading-agents/internal/orchestrator"
	"github.com/sousandrei/trading-agents/internal/tools/finnhub"
	"github.com/sousandrei/trading-agents/internal/tools/gemini"
	"github.com/sousandrei/trading-agents/internal/tools/simfin"
)

func main() {
	ctx := context.Background()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Println("failed to parse env vars: ", slog.Any("err", err))
		return
	}

	slog.Info("Starting Trading Agents Orchestrator")

	finnhubClient := finnhub.New(cfg.FinnhubAPIKey)
	simfinClient := simfin.New(cfg.SimfinAPIKey)

	geminiClient, err := gemini.New(
		ctx,
		finnhubClient,
		simfinClient,
		cfg.LLMModel,
	)
	if err != nil {
		log.Println("failed to create Gemini client: ", slog.Any("err", err))
		return
	}

	res, err := orchestrator.Analyze(ctx, geminiClient, "NVDA")
	if err != nil {
		log.Println("failed to analyze: ", slog.Any("err", err))
		return
	}

	fmt.Println("Finall Result:", res)
}
