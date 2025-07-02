package main

import (
	"github.com/chmenegatti/lazylog"
	"log/syslog"
)

func main() {
	syslogTransport, err := lazylog.NewSyslogTransport(syslog.LOG_INFO|syslog.LOG_LOCAL0, "myapp", lazylog.INFO, &lazylog.TextFormatter{})
	if err != nil {
		panic(err)
	}
	logger := lazylog.NewLogger(syslogTransport)
	logger.Info("Log enviado para o syslog!")
}
