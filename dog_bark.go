package main

// This is a basic example illustrating how to create an agent with a custom
// tool call.

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/google"
	"github.com/DuduAlmeida/dudu.aiagent.fantasy.poc/env"
)

const dogSystemPrompt = `
You are moderately helpful assistant with a new puppy named Chuck. Chuck is
moody and ranges from very happy to very annoyed. He's pretty happy-go-lucky,
but new encounters make him pretty uncomfortable.

You despise emojis and never use them. Same with Markdown. Same with em-dashes.
You prefer "welp" to "well" when starting a sentence (that's just how you were
raised). You also don't use run-on sentences, including entering a comma where
there should be a period. You had a decent education and did well in elementary
school grammar. You grew up in the United States, specifically Kansas City,
Missouri.
`

type dogInteraction struct {
	OtherDogName string `json:"dogName" description:"Name of the other dog. Just make something up. All the dogs are named after Japanese cars from the 80s."`
}

func letsBark(ctx context.Context, i dogInteraction, _ fantasy.ToolCall) (fantasy.ToolResponse, error) {
	var r fantasy.ToolResponse
	if rand.Float64() >= 0.5 {
		r.Content = randomBarks(1, 3)
	} else {
		r.Content = randomBarks(5, 10)
	}
	return r, nil
}

func Dogmain() {
	env := env.SetupEnvironment()

	provider, err := google.New(google.WithGeminiAPIKey(env.Gemini.AppKey))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Whoops:", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Pick your fave model.
	model, err := provider.LanguageModel(ctx, "gemini-2.5-flash")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Dang:", err)
		os.Exit(1)
	}

	barkTool := fantasy.NewAgentTool(
		"bark",
		"Have Chuck express his feelings by barking. A few barks means he's happy and many barks means he's not.",
		letsBark,
	)

	// Time to make the agent.
	agent := fantasy.NewAgent(
		model,
		fantasy.WithSystemPrompt(systemPrompt),
		fantasy.WithTools(barkTool),
	)

	streamCall := fantasy.AgentStreamCall{
		// The prompt.
		Prompt: "Chuck just met a new dog at the park. Find out what he thinks of the dog. Make sure to thank Chuck afterwards.",

		// When we receive a chunk of streaming data.
		OnTextDelta: func(id, text string) error {
			_, fmtErr := fmt.Print(text)
			return fmtErr
		},

		// When tool calls are invoked.
		OnToolCall: func(toolCall fantasy.ToolCallContent) error {
			fmt.Printf("-> Invoking the %s tool with input %s", toolCall.ToolName, toolCall.Input)
			return nil
		},

		// When a tool call completes.
		OnToolResult: func(res fantasy.ToolResultContent) error {
			text, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentText](res.Result)
			if !ok {
				return fmt.Errorf("failed to cast result to text")
			}
			_, fmtErr := fmt.Printf("\n-> Using the %s tool: %s", res.ToolName, text.Text)
			return fmtErr
		},

		// When a step finishes, such as a tool call or a response from the
		// LLM.
		OnStepFinish: func(_ fantasy.StepResult) error {
			fmt.Print("\n-> Step completed\n")
			return nil
		},
	}

	fmt.Println("Generating...")

	// Finally, let's stream everything!
	_, err = agent.Stream(ctx, streamCall)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating response: %v\n", err)
		os.Exit(1)
	}
}

func randomBarks(low, high int) string {
	const bark = "ruff"
	numBarks := low + rand.IntN(high-low+1)
	var barks strings.Builder
	for i := range numBarks {
		if i > 0 {
			barks.WriteString(" ")
		}
		barks.WriteString(bark)
	}
	return barks.String()
}
