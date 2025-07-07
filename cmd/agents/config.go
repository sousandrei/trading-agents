package main

import (
	"time"
)

type Config struct {
	Port      int    `env:"PORT" envDefault:"3001"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"text"`

	LLMProvider string `env:"LLM_PROVIDER" envDefault:"gemini"`
	LLMModel    string `env:"LLM_MODEL" envDefault:"gemini-2.5-flash"`

	// Ollama specific
	LLMApiUrl string `env:"LLM_API_URL" envDefault:"http://localhost:11434"`

	FinnhubAPIKey string `env:"FINNHUB_API_KEY"`
	SimfinAPIKey  string `env:"SIMFIN_API_KEY"`

	APITimeout  time.Duration `env:"API_TIMEOUT" envDefault:"10s"`
	APICacheTTL time.Duration `env:"API_CACHE_TTL" envDefault:"8760h"`
}
