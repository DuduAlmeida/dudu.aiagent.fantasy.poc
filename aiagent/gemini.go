package aiagent

import (
	"context"
	"fmt"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/google"
)

type GeminiAIAgent struct {
	providerID     AIProvider
	provider       fantasy.Provider
	defaultModelID string
	model          fantasy.LanguageModel
	appKey         string

	agent fantasy.Agent
}

func NewGeminiAIAgent(geminiAppKey string) AIAgent {
	return &GeminiAIAgent{
		providerID:     AIProviderGemini,
		appKey:         geminiAppKey,
		defaultModelID: "gemini-2.5-flash",
	}
}

func (g *GeminiAIAgent) Provider() AIProvider {
	return g.providerID
}

func (g *GeminiAIAgent) Bootstrap(ctx context.Context, opts ...fantasy.AgentOption) error {
	provider, err := google.New(google.WithGeminiAPIKey(g.appKey))
	if err != nil {
		return fmt.Errorf("error starting gemini provider: %s", err.Error())
	}

	model, err := provider.LanguageModel(ctx, g.defaultModelID)
	if err != nil {
		return fmt.Errorf("error starting gemini model %s: %s", g.defaultModelID, err.Error())
	}

	agent := fantasy.NewAgent(
		model,
		opts...,
	)

	g.provider = provider
	g.model = model
	g.agent = agent

	return nil
}

func (g *GeminiAIAgent) SubmitPrompt(ctx context.Context, prompt string) (err error) {
	streamCall := fantasy.AgentStreamCall{
		Prompt: prompt,

		OnTextDelta: func(id, text string) error {
			_, fmtErr := fmt.Print(text)
			return fmtErr
		},

		OnToolCall: func(toolCall fantasy.ToolCallContent) error {
			fmt.Printf("-> Invoking the %s tool with input %s", toolCall.ToolName, toolCall.Input)
			return nil
		},

		OnToolResult: func(res fantasy.ToolResultContent) error {
			text, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentText](res.Result)
			if !ok {
				return fmt.Errorf("failed to cast result to text")
			}
			_, fmtErr := fmt.Printf("\n-> Using the %s tool: %s", res.ToolName, text.Text)
			return fmtErr
		},

		OnStepFinish: func(_ fantasy.StepResult) error {
			fmt.Print("\n-> Step completed\n")
			return nil
		},
	}

	fmt.Println("Generating...")

	_, err = g.agent.Stream(ctx, streamCall)
	if err != nil {
		return fmt.Errorf("Error generating response: %v\n", err.Error())
	}

	return nil
}
