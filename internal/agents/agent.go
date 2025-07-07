package agents

import (
	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

type Agent struct {
	Prompt string

	Tools  bool
	Search bool

	Messages []llms.Message
}
