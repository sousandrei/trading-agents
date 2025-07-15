package llms

import (
	"context"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleModel Role = "model"
)

type Message struct {
	Role Role
	Text string
}

type Client interface {
	Generate(ctx context.Context, prompt string, opts ...GenerateOption) ([]Message, error)
}

type GenerateOptions struct {
	DryRun bool

	Search bool
	Tools  bool

	Messages []Message
}

type GenerateOption func(*GenerateOptions)

func WithDryRun() func(*GenerateOptions) {
	return func(c *GenerateOptions) {
		c.DryRun = true
	}
}

func WithSearch() func(*GenerateOptions) {
	return func(c *GenerateOptions) {
		c.Search = true
	}
}

func WithTools() func(*GenerateOptions) {
	return func(c *GenerateOptions) {
		c.Tools = true
	}
}

func WithMessages(messages []Message) func(*GenerateOptions) {
	return func(c *GenerateOptions) {
		c.Messages = messages
	}
}
