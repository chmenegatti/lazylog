package main

import (
	"log"
	"os"

	"github.com/chmenegatti/lazylog" // Importa a biblioteca
)

func main() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileLogger := lazylog.New()
	fileLogger.Out = file

	fileLogger.Info("Servidor iniciado.")
	fileLogger.Warn("A conexão com o banco de dados está lenta.")
	fileLogger.Error("Falha ao processar a requisição #123.")
	fileLogger.Info("Servidor finalizado.")

	log.Println("Logs foram escritos em app.log")
}
