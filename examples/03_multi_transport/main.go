package main

import (
	"fmt"
	"os"

	"github.com/chmenegatti/lazylog"
)

func main() {
	logger := lazylog.NewLogger()

	// Transporte para console (stdout)
	console := &lazylog.ConsoleTransport{
		Level:     lazylog.DEBUG,
		Formatter: &lazylog.TextFormatter{},
		ToStdErr:  false,
	}
	logger.AddTransport(console)

	// Transporte para arquivo
	fileTransport, err := lazylog.NewFileTransport("app.log", lazylog.INFO, &lazylog.JSONFormatter{})
	if err != nil {
		fmt.Println("Erro ao criar transporte de arquivo:", err)
		os.Exit(1)
	}
	defer fileTransport.Close()
	logger.AddTransport(fileTransport)

	logger.Debug("Mensagem de debug (apenas console)")
	logger.Info("Mensagem de info (console e arquivo)")
	logger.Warn("Mensagem de aviso")
	logger.Error("Mensagem de erro")

	// Log com metadata/contexto extra
	logger.ComFields(map[string]interface{}{
		"user":       "cesar",
		"request_id": 12345,
	}).Info("Log com contexto extra")
}
