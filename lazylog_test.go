package lazylog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/chmenegatti/lazylog"
)

func TestConsoleTransport_TextFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.TextFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	logger.Info("Hello World!")
	out := buf.String()
	if !strings.Contains(out, "Hello World!") || !strings.Contains(out, "[INFO]") {
		t.Errorf("output missing expected content: %s", out)
	}
}

func TestFileTransport_JSONFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.DEBUG,
		Formatter: &lazylog.JSONFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	logger.Debug("json test")
	var m map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &m)
	if err != nil || m["message"] != "json test" || m["level"] != "DEBUG" {
		t.Errorf("json output invalid: %v, %v", err, m)
	}
}

func TestComFieldsAndHooks(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.JSONFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	logger.AddHook(func(e *lazylog.Entry) {
		if e.Fields == nil {
			e.Fields = make(map[string]interface{})
		}
		e.Fields["hooked"] = true
	}, true)
	logger.ComFields(map[string]interface{}{"user": "cesar"}).Info("with fields")
	var m map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &m)
	if m["user"] != "cesar" || m["hooked"] != true {
		t.Errorf("fields or hook missing: %v", m)
	}
}

func TestCustomLevel(t *testing.T) {
	lazylog.RegisterLevel("NOTICE", 10)
	if lazylog.ParseLevel("NOTICE") != 10 {
		t.Error("custom level not registered")
	}
}

func TestTransportFilter(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.TextFormatter{},
	}
	filter := func(e *lazylog.Entry) bool {
		return strings.Contains(e.Message, "importante")
	}
	logger := lazylog.NewLogger(&lazylog.TransportWithFilter{Transport: tr, Filter: filter})
	logger.Info("não vai logar")
	logger.Info("log importante!")
	out := buf.String()
	if !strings.Contains(out, "importante") || strings.Contains(out, "não vai logar") {
		t.Errorf("filter failed: %s", out)
	}
}

func TestContextSupport(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.JSONFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	ctx := context.WithValue(context.Background(), "trace_id", "abc-123")
	logger.InfoCtx(ctx, "ctx test", nil)
	var m map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &m)
	if m["trace_id"] != "abc-123" {
		t.Errorf("context trace_id missing: %v", m)
	}
}

func TestWithFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.TextFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	logger.WithFormatter(&lazylog.TextFormatter{TimestampFormat: time.RFC822}).Info("custom format")
	out := buf.String()
	if !strings.Contains(out, "custom format") || !strings.Contains(out, "[INFO]") || !strings.Contains(out, "Jul") {
		t.Errorf("custom formatter not applied: %s", out)
	}
}

func TestNestedFieldsJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := &lazylog.WriterTransport{
		Writer:    buf,
		Level:     lazylog.INFO,
		Formatter: &lazylog.JSONFormatter{},
	}
	logger := lazylog.NewLogger(tr)
	logger.ComFields(map[string]interface{}{
		"user": "cesar",
		"request": map[string]interface{}{
			"id": 123,
			"ip": "1.2.3.4",
		},
	}).Info("nested fields")
	var m map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &m)
	if req, ok := m["request"].(map[string]interface{}); !ok || req["id"] != float64(123) {
		t.Errorf("nested fields not serialized: %v", m)
	}
}
