package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chmenegatti/lazylog"
)

func main() {
	// Cria transportes para console e arquivo
	console := &lazylog.ConsoleTransport{
		Level:     lazylog.DEBUG,
		Formatter: &lazylog.TextFormatter{},
	}
	fileTransport, err := lazylog.NewFileTransport("app.log", lazylog.INFO, &lazylog.JSONFormatter{})
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de log:", err)
		os.Exit(1)
	}
	defer fileTransport.Close()

	logger := lazylog.NewLogger(console, fileTransport)

	// Hook para adicionar um campo de tempo de execução
	logger.AddHook(func(e *lazylog.Entry) {
		if e.Fields == nil {
			e.Fields = make(map[string]any)
		}
		e.Fields["runtime"] = time.Now().UnixNano()
	}, true) // before

	// Hook para imprimir no console toda vez que um erro for logado
	logger.AddHook(func(e *lazylog.Entry) {
		if e.Level == lazylog.ERROR {
			fmt.Println("[ALERTA] Um erro foi registrado:", e.Message)
		}
	}, false) // after

	logger.Info("Log normal com hooks")
	logger.Error("Erro de exemplo com hooks")
	logger.ComFields(map[string]any{"user": "cesar"}).Warn("Log com contexto e hooks")
}
