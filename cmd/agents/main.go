package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/caarlos0/env"
	"github.com/sousandrei/trading-agents/internal/orchestrator"
	"github.com/sousandrei/trading-agents/internal/tools/apiclient"
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

	switch cfg.LogFormat {
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(log.Writer(), nil)))
	case "text":
		slog.SetDefault(slog.New(slog.NewTextHandler(log.Writer(), nil)))
	default:
		log.Println("invalid log format, must be 'json' or 'text'")
		return
	}

	slog.Info("Starting Trading Agents Orchestrator")

	finnhubClient, err := finnhub.New(cfg.FinnhubAPIKey, apiclient.WithCache("data/finnhub.bin"))
	if err != nil {
		log.Println("failed to create Finnhub client: ", slog.Any("err", err))
		return
	}

	simfinClient, err := simfin.New(cfg.SimfinAPIKey, apiclient.WithCache("data/simfin.bin"))
	if err != nil {
		log.Println("failed to create Simfin client: ", slog.Any("err", err))
		return
	}

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
