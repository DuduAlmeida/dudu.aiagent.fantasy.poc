package env

type Enviroment struct {
	Gemini GeminiConfig `toml:"gemini"`
}

// Estrutura para a tabela [database]
type GeminiConfig struct {
	AppKey string `toml:"app-key"`
}
