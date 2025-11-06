package main

import (
	"context"

	"github.com/DuduAlmeida/dudu.aiagent.fantasy.poc/aiagent"
	"github.com/DuduAlmeida/dudu.aiagent.fantasy.poc/env"
	"github.com/DuduAlmeida/dudu.aiagent.fantasy.poc/experiments"
)

func main() {
	ctx := context.Background()

	env := env.SetupEnvironment()
	agent := aiagent.NewGeminiAIAgent(env.Gemini.AppKey)

	experiments.SchedullerExperiment(ctx, agent)
}
