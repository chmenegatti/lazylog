package lazylog_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/chmenegatti/lazylog"
)

func BenchmarkLogger_Info(b *testing.B) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.TextFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("mensagem de benchmark")
	}
}

func BenchmarkLogger_JSONFields(b *testing.B) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.JSONFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	fields := map[string]any{
		"user":    "cesar",
		"request": map[string]any{"id": 123, "ip": "1.2.3.4"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.ComFields(fields).Info("mensagem estruturada")
	}
}

func BenchmarkLogger_Parallel(b *testing.B) {
	tr := &lazylog.WriterTransport{
		Writer:    io.Discard,
		Level:     lazylog.INFO,
		Formatter: &lazylog.TextFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("log paralelo")
		}
	})
}

func BenchmarkLogger_WithHooks(b *testing.B) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.TextFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	logger.AddHook(func(e *lazylog.Entry) {
		e.Fields = map[string]interface{}{"bench": true}
	}, true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("mensagem com hook")
	}
}

func BenchmarkLogger_WithFormatter(b *testing.B) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.TextFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	formatter := &lazylog.TextFormatter{TimestampFormat: time.RFC822}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFormatter(formatter).Info("mensagem customizada")
	}
}
