package aiagent

import (
	"context"

	"charm.land/fantasy"
)

// In this POC I left some points to improve, like create an interface that fantasy's contracts implements, enabling use others libs than fantasy
type AIAgent interface {
	Provider() AIProvider
	Bootstrap(ctx context.Context, opts ...fantasy.AgentOption) error
	SubmitPrompt(ctx context.Context, prompt string) (*fantasy.AgentResult, error)
}

type AIProvider string

func (a AIProvider) String() string {
	return string(a)
}

const (
	AIProviderGemini AIProvider = "gemini"
)
