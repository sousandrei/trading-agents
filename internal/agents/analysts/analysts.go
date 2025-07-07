package analysts

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/sousandrei/trading-agents/internal/agents"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

func Run(
	ctx context.Context,
	llm llms.Client,
	ticker string,
	opts ...llms.GenerateOption,
) (map[string]agents.Agent, error) {
	if opts == nil {
		opts = []llms.GenerateOption{}
	}

	g, ctx := errgroup.WithContext(ctx)

	analysts := map[string]agents.Agent{
		"fundamentals": {
			Prompt:   fundamentalsPrompt,
			Tools:    true,
			Messages: []llms.Message{},
		},
		"market": {
			Prompt:   marketPrompt,
			Tools:    true,
			Messages: []llms.Message{},
		},
		"social_media": {
			Prompt:   socialMediaPrompt,
			Search:   true,
			Messages: []llms.Message{},
		},
		"news": {
			Prompt:   newsPrompt,
			Search:   true,
			Messages: []llms.Message{},
		},
	}

	for name, agent := range analysts {
		g.Go(func() error {
			if agent.Tools {
				opts = append(opts, llms.WithTools())
			}

			if agent.Search {
				opts = append(opts, llms.WithSearch())
			}

			res, err := llm.Generate(ctx, agent.Prompt, opts...)
			if err != nil {
				return fmt.Errorf("error running analyst %s: %w", agent.Prompt, err)
			}

			prompt := fmt.Sprintf("%s\nStock in question: %s", agent.Prompt, ticker)

			analysts[name] = agents.Agent{
				Prompt:   prompt,
				Tools:    agent.Tools,
				Search:   agent.Search,
				Messages: res,
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to run analysts: %w", err)
	}

	// TODO: remove, dev only
	for name, agent := range analysts {
		agents.WriteMessagesToFile("analysts", name, agent.Messages)
	}

	return analysts, nil
}

func AppendOutput(prompt string, analysts map[string]agents.Agent) string {
	for name, agent := range analysts {
		for _, message := range agent.Messages {
			if message.Role != llms.RoleModel {
				continue
			}

			prompt += fmt.Sprintf("\n\n%s analyst: %s", name, message.Text)
		}
	}

	return prompt
}
