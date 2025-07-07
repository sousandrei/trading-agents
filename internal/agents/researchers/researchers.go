package researchers

import (
	"context"
	"fmt"

	"github.com/sousandrei/trading-agents/internal/agents"
	"github.com/sousandrei/trading-agents/internal/agents/analysts"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

func Run(
	ctx context.Context,
	llm llms.Client,
	ticker string,
	analystAgents map[string]agents.Agent,
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
				prompt = analysts.AppendOutput(agent.Prompt, analystAgents)
			} else {
				for r, a := range researchers {
					if r == researcher {
						continue
					}

					prompt = fmt.Sprintf("last %s argument: %s", r, a.Messages[len(a.Messages)-1].Text)
				}
			}

			opts = append(opts, llms.WithMessages(agent.Messages))

			res, err := llm.Generate(ctx, prompt, opts...)
			if err != nil {
				return nil, fmt.Errorf("error running bull researcher: %w", err)
			}

			researchers[researcher] = agents.Agent{
				Prompt:   researchers[researcher].Prompt,
				Messages: res,
			}
		}
	}

	res, err := llm.Generate(ctx, AppendOutput(managerPrompt, researchers), opts...)
	if err != nil {
		return nil, fmt.Errorf("error running manager: %w", err)
	}

	researchers["manager"] = agents.Agent{
		Prompt:   managerPrompt,
		Messages: res,
	}

	return researchers, nil
}

func AppendOutput(prompt string, researchers map[string]agents.Agent) string {
	prompt += fmt.Sprintf("%s\n\n\n\nResearch team debate:", prompt)

	for round := range 3 {
		for name, agent := range researchers {
			prompt += fmt.Sprintf("\n\n%s: %s", name, agent.Messages[round].Text)
		}
	}

	return prompt
}

func AppendManagerOutput(prompt string, researchers map[string]agents.Agent) string {
	for name, agent := range researchers {
		if name == "manager" {
			prompt += fmt.Sprintf("\n\nresearch manager report: %s", agent.Messages[len(agent.Messages)-1].Text)
			break
		}
	}

	return prompt
}
