package main

import (
	"log/syslog"

	"github.com/chmenegatti/lazylog"
)

func main() {
	syslogTransport, err := lazylog.NewSyslogTransport(syslog.LOG_INFO|syslog.LOG_LOCAL0, "myapp", lazylog.INFO, &lazylog.TextFormatter{})
	if err != nil {
		panic(err)
	}
	logger := lazylog.NewLogger(syslogTransport)
	logger.Info("Log enviado para o syslog!")
}
