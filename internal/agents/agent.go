package agents

import (
	"fmt"
	"os"

	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

type Agent struct {
	Prompt string

	Tools  bool
	Search bool

	Messages []llms.Message
}

func WriteMessagesToFile(team, name string, messages []llms.Message) error {
	history := ""

	for i, message := range messages {
		if i > 0 {
			history += "\n===SEP===\n"
		}
		history += fmt.Sprintf("%s\n---\n", message.Role)
		history += message.Text
	}

	return os.WriteFile(fmt.Sprintf("data/prompts/%s/%s.txt", team, name), []byte(history), 0644)
}
