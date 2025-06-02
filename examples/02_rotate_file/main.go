package main

import (
	"log"
	"strconv"
	"time"

	"github.com/chmenegatti/lazylog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./meu_app_rotacionado.log",
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	logger := lazylog.New()
	logger.Out = lumberjackLogger
	// Usando o formato JSON para variar
	logger.Formatter = &lazylog.JSONFormatter{}

	log.Println("Iniciando a geração de logs...")
	for i := 0; i < 50; i++ {
		logger.Info("Esta é a mensagem de log de número " + strconv.Itoa(i+1))
		time.Sleep(20 * time.Millisecond)
	}

	log.Println("Logs rotacionados foram escritos em meu_app_rotacionado.log")
}
