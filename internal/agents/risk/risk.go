package risk

import (
	"context"
	"fmt"

	"github.com/sousandrei/trading-agents/internal/agents"
	"github.com/sousandrei/trading-agents/internal/agents/analysts"
	"github.com/sousandrei/trading-agents/internal/agents/trader"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

func Run(
	ctx context.Context,
	llm llms.Client,
	ticker string,
	analystAgents map[string]agents.Agent,
	traderAgent *agents.Agent,
	opts ...llms.GenerateOption,
) (map[string]agents.Agent, error) {
	if opts == nil {
		opts = []llms.GenerateOption{}
	}

	riskAgents := map[string]agents.Agent{
		"aggressive": {
			Prompt:   aggressivePrompt,
			Messages: []llms.Message{},
		},
		"conservative": {
			Prompt:   conservativePrompt,
			Messages: []llms.Message{},
		},
		"neutral": {
			Prompt:   neutralPrompt,
			Messages: []llms.Message{},
		},
	}

	names := []string{"aggressive", "conservative", "neutral"}

	for round := range 3 {
		for _, name := range names {
			agent := riskAgents[name]

			var prompt string

			if round == 0 {
				prompt = fmt.Sprintf("%s\nStock in question: %s", agent.Prompt, ticker)
				prompt = analysts.AppendOutput(prompt, analystAgents)
				prompt = trader.AppendOutput(prompt, traderAgent)
			} else {
				for _, a := range names {
					if a == name {
						continue
					}

					prompt = fmt.Sprintf("last %s argument: %s", a, "SUP")
				}
			}

			opts = append(opts, llms.WithMessages(agent.Messages))

			res, err := llm.Generate(ctx, prompt, opts...)
			if err != nil {
				return nil, fmt.Errorf("error running bull researcher: %w", err)
			}

			riskAgents[name] = agents.Agent{
				Prompt:   riskAgents[name].Prompt,
				Messages: res,
			}
		}
	}

	res, err := llm.Generate(ctx, AppendOutput(managerPrompt, riskAgents), opts...)
	if err != nil {
		return nil, fmt.Errorf("error running manager: %w", err)
	}

	riskAgents["manager"] = agents.Agent{
		Prompt:   managerPrompt,
		Messages: res,
	}

	return riskAgents, nil
}

func AppendOutput(prompt string, risk map[string]agents.Agent) string {
	prompt += fmt.Sprintf("%s\n\n\n\nRisk team debate:", prompt)

	for round := range 3 {
		for name, agent := range risk {
			prompt += fmt.Sprintf("\n\n%s: %s", name, agent.Messages[round].Text)
		}
	}

	return prompt
}
