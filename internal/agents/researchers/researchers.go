package researchers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sousandrei/trading-agents/internal/agents"
	"github.com/sousandrei/trading-agents/internal/agents/analysts"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
	"github.com/sousandrei/trading-agents/internal/types"
)

func Run(
	ctx context.Context,
	llm llms.Client,
	analystAgents map[string]agents.Agent,
	position types.Position,
	opts ...llms.GenerateOption,
) (map[string]agents.Agent, error) {
	researchers := map[string]agents.Agent{
		"bull": {
			Prompt:   bullPrompt,
			Messages: []llms.Message{},
		},
		"bear": {
			Prompt:   bearPrompt,
			Messages: []llms.Message{},
		},
	}

	if opts == nil {
		opts = []llms.GenerateOption{}
	}

	for round := range 3 {
		for _, researcher := range []string{"bull", "bear"} {
			agent := researchers[researcher]

			var prompt string

			if round == 0 {
				prompt = fmt.Sprintf("%s\n%s", agent.Prompt, position)
				prompt = analysts.AppendOutput(prompt, analystAgents)
			} else {
				for r, a := range researchers {
					if r == researcher {
						continue
					}

					prompt = fmt.Sprintf("last %s argument: %s", r, a.Messages[len(a.Messages)-1].Text)
				}
			}

			opts := append(opts, llms.WithMessages(agent.Messages))

			slog.Info("Running researcher", "name", researcher, "round", round, "ticker", position.Ticker)

			res, err := llm.Generate(ctx, prompt, opts...)
			if err != nil {
				return nil, fmt.Errorf("error running %s researcher: %w", researcher, err)
			}

			researchers[researcher] = agents.Agent{
				Prompt:   researchers[researcher].Prompt,
				Messages: res,
			}
		}
	}

	res, err := llm.Generate(ctx, AppendOutput(managerPrompt, researchers), opts...)
	if err != nil {
		return nil, fmt.Errorf("error running research manager: %w", err)
	}

	researchers["manager"] = agents.Agent{
		Prompt:   managerPrompt,
		Messages: res,
	}

	// TODO: remove, dev only
	for name, agent := range researchers {
		agents.WriteMessagesToFile("researchers", name, agent.Messages)
	}

	return researchers, nil
}

func AppendOutput(prompt string, researchers map[string]agents.Agent) string {
	prompt += fmt.Sprintf("%s\n\n\n\nResearch team debate:", prompt)

	modelMessages := map[string][]llms.Message{}

	for name, agent := range researchers {
		modelMessages[name] = []llms.Message{}
		for _, message := range agent.Messages {
			if message.Role == llms.RoleModel {
				modelMessages[name] = append(modelMessages[name], message)
			}
		}
	}

	for round := range 3 {
		for _, name := range []string{"bull", "bear"} {
			prompt += fmt.Sprintf("\n\n#### %s Analyst Argument (Round %d):\n%s", name, round+1, modelMessages[name][round].Text)
		}
	}

	return prompt
}

func AppendManagerOutput(prompt string, researchers map[string]agents.Agent) string {
	for name, agent := range researchers {
		if name == "manager" {
			prompt += fmt.Sprintf("\n\n### Research Manager Report:\n%s", agent.Messages[len(agent.Messages)-1].Text)
			break
		}
	}

	return prompt
}
