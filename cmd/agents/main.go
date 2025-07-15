package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/caarlos0/env/v11"

	"github.com/sousandrei/trading-agents/internal/orchestrator"
	"github.com/sousandrei/trading-agents/internal/server"
	"github.com/sousandrei/trading-agents/internal/tools/apiclient"
	"github.com/sousandrei/trading-agents/internal/tools/finnhub"
	"github.com/sousandrei/trading-agents/internal/tools/gemini"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
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

	llmClient, err := createLLMClient(ctx, cfg)
	if err != nil {
		log.Println("failed to create LLM client: ", slog.Any("err", err))
		return
	}

	o := orchestrator.New(llmClient)

	s := server.New(ctx, cfg.Server, o.Handler)
	go func() {
		slog.Info("Starting server", slog.String("addr", s.Addr))
		if err := s.ListenAndServe(); err != nil {
			log.Println("failed to start server: ", slog.Any("err", err))
			return
		}
		slog.Info("Server stopped")
	}()

	<-ctx.Done()

}

func createLLMClient(ctx context.Context, cfg Config) (llms.Client, error) {
	finnhubClient, err := finnhub.New(cfg.FinnhubAPIKey, apiclient.WithCache("data/finnhub.bin"))
	if err != nil {
		return nil, fmt.Errorf("failed to create Finnhub client: %w", err)
	}

	simfinClient, err := simfin.New(cfg.SimfinAPIKey, apiclient.WithCache("data/simfin.bin"))
	if err != nil {
		return nil, fmt.Errorf("failed to create Simfin client: %w", err)
	}

	switch cfg.LLMProvider {
	case "gemini":
		return gemini.New(
			ctx,
			finnhubClient,
			simfinClient,
			cfg.LLMModel,
		)
	case "ollama":
		return nil, fmt.Errorf("support for Ollama is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.LLMProvider)
	}
}
