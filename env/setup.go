package env

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func SetupEnvironment() Enviroment {
	var env Enviroment

	_, err := toml.DecodeFile("config.toml", &env)
	if err != nil {
		panic(fmt.Sprintf("Erro ao ler o arquivo de configuração TOML: %v", err))
	}

	return env
}
