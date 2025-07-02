package main

import (
	"fmt"
	"os"

	"github.com/chmenegatti/lazylog"
)

func main() {
	fileTransport, err := lazylog.NewFileTransport("app.log", lazylog.INFO, &lazylog.TextFormatter{})
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de log:", err)
		os.Exit(1)
	}
	defer fileTransport.Close()

	logger := lazylog.NewLogger(fileTransport)

	logger.Info("Servidor iniciado.")
	logger.Warn("A conexão com o banco de dados está lenta.")
	logger.Error("Falha ao processar a requisição #123.")
	logger.Info("Servidor finalizado.")
}
