package gemini

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/genai"

	"github.com/sousandrei/trading-agents/internal/tools/apiclient"
	"github.com/sousandrei/trading-agents/internal/tools/finnhub"
	"github.com/sousandrei/trading-agents/internal/tools/llms"
	"github.com/sousandrei/trading-agents/internal/tools/simfin"
)

type Client struct {
	gemini *genai.Client

	finnhub *finnhub.Client
	simfin  *simfin.Client

	model string
}

func New(
	ctx context.Context,
	finnhub *finnhub.Client,
	simfin *simfin.Client,
	model string,
) (*Client, error) {
	c, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			Backend: genai.BackendVertexAI,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		gemini:  c,
		finnhub: finnhub,
		simfin:  simfin,
		model:   model,
	}, nil
}

func (c *Client) Generate(
	ctx context.Context,
	prompt string,
	opts ...llms.GenerateOption,
) ([]llms.Message, error) {
	options := &llms.GenerateOptions{}

	for _, opt := range opts {
		opt(options)
	}

	messages := []llms.Message{}
	content := []*genai.Content{}

	if options.Messages != nil {
		for _, msg := range options.Messages {
			messages = append(messages, msg)

			content = append(content, &genai.Content{
				Role:  string(msg.Role),
				Parts: []*genai.Part{{Text: msg.Text}},
			})
		}
	}

	messages = append(messages, llms.Message{
		Role: llms.RoleUser,
		Text: prompt,
	})

	content = append(content, &genai.Content{
		Role:  genai.RoleUser,
		Parts: []*genai.Part{{Text: prompt}},
	})

	if options.DryRun {
		newMsg := llms.Message{
			Role: llms.RoleModel,
			Text: "Dry run mode enabled. No actual generation will occur.",
		}

		messages = append(messages, newMsg)
		return messages, nil

	}

	var fns map[string]apiclient.Fetch
	var tools []*genai.Tool

	if options.Search && options.Tools {
		return nil, fmt.Errorf("cannot use both search and tools together")
	}

	if options.Search {
		tools = searchTool()
	} else if options.Tools {
		fns, tools = c.apiTools()
	}

	res, err := c.gemini.Models.GenerateContent(
		ctx,
		c.model,
		content,
		&genai.GenerateContentConfig{
			Tools: tools,
		},
	)
	if err != nil {
		return nil, err
	}

	for {
		slog.Info("Received response", "text", res.Text(), "function_calls", res.FunctionCalls())

		if len(res.FunctionCalls()) == 0 {
			break
		}

		parts := []*genai.Part{
			{
				Text: prompt,
			},
		}

		for _, call := range res.FunctionCalls() {
			fn, ok := fns[call.Name]
			if !ok {
				return nil, fmt.Errorf("function %s not found", call.Name)
			}

			slog.Info("Calling function", "name", call.Name, "args", call.Args)

			res, err := fn(ctx, call.Args)
			if err != nil {
				return nil, fmt.Errorf("error calling function %s: %w", call.Name, err)
			}

			parts = append(parts, &genai.Part{
				FunctionResponse: &genai.FunctionResponse{
					ID:       call.ID,
					Name:     call.Name,
					Response: res,
				},
			})
		}

		content = []*genai.Content{
			{
				Role:  genai.RoleUser,
				Parts: parts,
			},
		}

		res, err = c.gemini.Models.GenerateContent(
			ctx,
			c.model,
			content,
			&genai.GenerateContentConfig{
				Tools: tools,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error generating content: %w", err)
		}
	}

	messages = append(messages, llms.Message{
		Role: llms.RoleModel,
		Text: res.Text(),
	})

	return messages, nil
}
