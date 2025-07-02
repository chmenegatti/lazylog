package main

import (
	"os"

	"github.com/chmenegatti/lazylog"
)

func main() {
	// Exemplo: configuração via arquivo JSON
	cfg, err := lazylog.LoadLoggerConfigJSON("logger_config.json")
	if err != nil {
		panic(err)
	}
	logger, err := lazylog.NewLoggerFromConfig(cfg)
	if err != nil {
		panic(err)
	}
	logger.Info("Logger configurado via JSON!")

	// Exemplo: configuração via arquivo YAML
	cfgYAML, err := lazylog.LoadLoggerConfigYAML("logger_config.yaml")
	if err == nil {
		loggerYAML, _ := lazylog.NewLoggerFromConfig(cfgYAML)
		loggerYAML.Info("Logger configurado via YAML!")
	}

	// Exemplo: uso de child logger
	child := logger.WithFields(map[string]any{"service": "auth", "env": os.Getenv("ENV")})
	child.Info("Log do serviço de autenticação")
	child.Error("Erro no serviço de autenticação", map[string]any{"code": 401})
}
