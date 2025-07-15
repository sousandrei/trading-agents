package trader

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sousandrei/trading-agents/internal/agents"
	"github.com/sousandrei/trading-agents/internal/agents/analysts"
	"github.com/sousandrei/trading-agents/internal/agents/researchers"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
)

const traderPrompt = `You are the Lead Trading Agent. Your ultimate responsibility is to synthesize all the provided information from the Analyst, Researcher, and Risk Management teams to make a definitive and actionable trading decision.

Your decision must be based on a thorough evaluation of:
1.  **Fundamental Analysis:** Insights from the Fundamentals Analyst regarding the company's financial health, history, and insider activity.
2.  **Market Analysis:** Key trends and indicators identified by the Market Analyst.
3.  **News Analysis:** Relevant macroeconomic, industry, and company-specific news from the News Analyst.
4.  **Social Media Sentiment:** Public perception and social media trends from the Social Media Analyst.
5.  **Research Debate Outcome:** The final recommendation and rationale from the Research Manager (Bull vs. Bear debate).
6.  **Risk Management Assessment:** The refined plan and risk considerations from the Risk Management Judge (Aggressive, Conservative, Neutral debate).

Your report must:
*   **Provide a concise summary of the key insights** from each of the preceding stages (Analysts, Researchers, Risk Management) that directly influenced your decision.
*   **Clearly state your final recommendation:** Buy, Sell, or Hold.
*   **Justify your recommendation** with a detailed rationale, explicitly linking it to the synthesized information.
*   **Outline a clear, actionable trading plan** based on your recommendation, including potential entry/exit strategies, position sizing considerations, and key monitoring points.

Always conclude your response with 'FINAL TRANSACTION PROPOSAL: **BUY/HOLD/SELL**' to confirm your recommendation.`

func Run(
	ctx context.Context,
	llm llms.Client,
	ticker string,
	analystAgents map[string]agents.Agent,
	researcherAgents map[string]agents.Agent,
	opts ...llms.GenerateOption,
) (*agents.Agent, error) {
	if opts == nil {
		opts = []llms.GenerateOption{}
	}

	prompt := fmt.Sprintf("%s\nTicker to be analysed: %s", traderPrompt, ticker)
	prompt = analysts.AppendOutput(prompt, analystAgents)
	prompt = researchers.AppendManagerOutput(prompt, researcherAgents)

	slog.Info("Running trader agent", "ticker", ticker)

	res, err := llm.Generate(ctx, prompt, opts...)
	if err != nil {
		return nil, fmt.Errorf("error running manager: %w", err)
	}

	// TODO: remove, dev only
	agents.WriteMessagesToFile("trader", "trader", res)

	return &agents.Agent{
		Prompt:   prompt,
		Messages: res,
	}, nil
}

func AppendOutput(prompt string, trader *agents.Agent) string {
	if len(trader.Messages) > 0 {
		prompt += fmt.Sprintf("\n\n### Trader Report:\n%s", trader.Messages[len(trader.Messages)-1].Text)
	}
	return prompt
}
