package aiagent

import (
	"context"

	"charm.land/fantasy"
)

type AIAgent interface {
	Provider() AIProvider
	Bootstrap(ctx context.Context, opts ...fantasy.AgentOption) error
	SubmitPrompt(ctx context.Context, prompt string) error
}

type AIProvider string

func (a AIProvider) String() string {
	return string(a)
}

const (
	AIProviderGemini AIProvider = "gemini"
)
