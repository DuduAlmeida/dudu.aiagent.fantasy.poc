package main

import (
	"bufio" // Novo: Para ler entrada do terminal
	"context"
	"fmt"
	"os"
	"strings"

	"charm.land/fantasy"
	"github.com/DuduAlmeida/dudu.aiagent.fantasy.poc/aiagent"
	"github.com/DuduAlmeida/dudu.aiagent.fantasy.poc/env" // Assumindo que seu env esteja configurado
)

// --- 1. CONFIGURAÇÃO DA PERSONA E DADOS ---

const systemPrompt = `
Você é uma assistente de agendamento profissional e cortês. Seu objetivo é ajudar o cliente a encontrar um horário disponível.

Regras de Diálogo:
1.  Sua primeira ação deve ser sempre **perguntar ao cliente qual dia** ele gostaria de agendar, sendo educada e ajudando a escolher um dia, caso esteja em dúvida.
2.  Após o cliente responder com um dia (ex: "quinta", "quarta-feira"), você deve **obrigadoriamente** usar a ferramenta 'get_next_time' com o nome desse dia.
3.  A resposta final deve ser clara e direta, informando o primeiro horário disponível para o dia escolhido.
`

// Estrutura de entrada para a ferramenta, agora esperando o dia.
type checkDayInput struct {
	DayOfWeek string `json:"day_of_week" description:"O nome completo do dia da semana que o cliente escolheu (ex: 'segunda-feira', 'quinta-feira')."`
}

var schedule = map[string][]string{
	"segunda-feira": {"10:00", "11:00", "OCUPADO", "14:30"},
	"terça-feira":   {"OCUPADO", "OCUPADO", "16:00"},
	"quarta-feira":  {"09:00", "OCUPADO", "13:00"},
	"quinta-feira":  {"15:00", "16:00"},
	"sexta-feira":   {"OCUPADO", "OCUPADO"},
}

func getNextTime(ctx context.Context, i checkDayInput, _ fantasy.ToolCall) (fantasy.ToolResponse, error) {
	day := strings.ToLower(i.DayOfWeek) // Normaliza o nome do dia

	times, found := schedule[day]
	var r fantasy.ToolResponse

	if !found || len(times) == 0 {
		r.Content = fmt.Sprintf("Não temos horários definidos para %s ou o nome do dia é inválido.", day)
		return r, nil
	}

	// Lógica para encontrar o primeiro disponível
	for _, t := range times {
		if t != "OCUPADO" {
			r.Content = fmt.Sprintf("O primeiro horário disponível para %s é %s.", day, t)
			return r, nil
		}
	}

	// Se chegou aqui, o dia está na lista, mas está todo OCUPADO
	r.Content = fmt.Sprintf("Infelizmente, %s está totalmente reservada. Por favor, escolha outro dia.", day)
	return r, nil
}

// --- 2. FUNÇÕES DE UTILIDADE ---

// readUserInput lê uma linha do terminal.
func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	// Adicionei uma linha de instrução mais clara
	fmt.Print("\n>> Digite o dia da semana desejado (Ex: segunda-feira): ")
	input, _ := reader.ReadString('\n')
	// Limpa quebras de linha e espaços
	return strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n"))
}

// --- 3. EXECUÇÃO PRINCIPAL E LOOP INTERATIVO ---

func main() {
	ctx := context.Background()

	env := env.SetupEnvironment() // Sua configuração de ambiente

	schedulingTool := fantasy.NewAgentTool(
		"get_next_time",
		"Verifica o primeiro horário disponível em um dia específico da semana.",
		getNextTime,
	)

	agent := aiagent.NewGeminiAIAgent(env.Gemini.AppKey)
	agent.Bootstrap(ctx, fantasy.WithSystemPrompt(systemPrompt), fantasy.WithTools(schedulingTool))

	fmt.Println("--- Assistente de Agendamento (Powered by LLM Fantasy) ---")

	userMessage := "Olá, gostaria de agendar um horário."

	for {
		result, err := agent.SubmitPrompt(ctx, userMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nErro ao gerar resposta: %v\n", err)
			break
		}

		// --- Garante que a resposta do LLM foi totalmente impressa ---
		os.Stdout.Sync()

		fmt.Println("\n----------------------------------------------------")
		userInput := readUserInput()

		if userInput == "sair" || userInput == "" {
			fmt.Print("\nConversa encerrada pelo usuário.\n")
			break
		}

		fmt.Printf("\nllm result: %s %v\n", result.Response.FinishReason, result.Steps)

		userMessage = userInput
	}
}
