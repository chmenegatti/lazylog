package lazylog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Hook func(entry *Entry)

// TransportErrorHook é chamado quando um transporte falha ao gravar.
type TransportErrorHook func(entry *Entry, transport Transport, err error)

// StacktraceConfig permite ativar stacktrace automático para níveis específicos.
type StacktraceConfig struct {
	Enabled bool
	Levels  map[Level]bool // Níveis que devem incluir stacktrace
}

// Logger é o logger principal, thread-safe para uso concorrente.
type Logger struct {
	mu          sync.RWMutex
	transports  []Transport
	beforeHooks []Hook
	afterHooks  []Hook
	errorHooks  []TransportErrorHook
	stacktrace  StacktraceConfig
}

// NewLogger cria um logger com zero ou mais transportes.
func NewLogger(transports ...Transport) *Logger {
	return &Logger{
		transports: transports,
	}
}

// EnableStacktrace ativa stacktrace automático para os níveis informados.
func (l *Logger) EnableStacktrace(levels ...Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.stacktrace.Enabled = true
	l.stacktrace.Levels = make(map[Level]bool)
	for _, lvl := range levels {
		l.stacktrace.Levels[lvl] = true
	}
}

// AddTransport adiciona um novo transporte ao logger.
func (l *Logger) AddTransport(t Transport) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.transports = append(l.transports, t)
}

// RemoveTransport remove um transporte do logger (por comparação de ponteiro).
func (l *Logger) RemoveTransport(t Transport) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i, tr := range l.transports {
		if tr == t {
			l.transports = append(l.transports[:i], l.transports[i+1:]...)
			return
		}
	}
}

// AddHook adiciona um hook para ser executado antes ou depois do log.
func (l *Logger) AddHook(hook Hook, before bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if before {
		l.beforeHooks = append(l.beforeHooks, hook)
	} else {
		l.afterHooks = append(l.afterHooks, hook)
	}
}

// AddErrorHook adiciona um hook para erros de transporte.
func (l *Logger) AddErrorHook(hook TransportErrorHook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errorHooks = append(l.errorHooks, hook)
}

// Close fecha todos os transportes que implementam io.Closer.
func (l *Logger) Close() error {
	l.mu.RLock()
	transports := make([]Transport, len(l.transports))
	copy(transports, l.transports)
	l.mu.RUnlock()

	var firstErr error
	for _, t := range transports {
		if closer, ok := t.(io.Closer); ok {
			if err := closer.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// snapshot retorna cópias locais dos campos protegidos para uso seguro fora do lock.
type logSnapshot struct {
	transports  []Transport
	beforeHooks []Hook
	afterHooks  []Hook
	errorHooks  []TransportErrorHook
	stacktrace  StacktraceConfig
}

func (l *Logger) snapshot() logSnapshot {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return logSnapshot{
		transports:  l.transports,
		beforeHooks: l.beforeHooks,
		afterHooks:  l.afterHooks,
		errorHooks:  l.errorHooks,
		stacktrace:  l.stacktrace,
	}
}

// dispatchEntry é a lógica centralizada de despacho de entry para transportes e hooks.
func dispatchEntry(snap logSnapshot, entry *Entry, formatter Formatter) {
	for _, hook := range snap.beforeHooks {
		hook(entry)
	}
	for _, t := range snap.transports {
		if entry.Level >= t.MinLevel() {
			var err error
			if formatter != nil {
				// Formata com o formatter customizado e escreve diretamente,
				// sem alterar o formatter do transporte (thread-safe).
				formatted, fmtErr := formatter.Format(entry)
				if fmtErr != nil {
					err = fmtErr
				} else {
					err = writeFormatted(t, formatted)
				}
			} else {
				err = t.WriteLog(entry)
			}
			if err != nil {
				for _, eh := range snap.errorHooks {
					eh(entry, t, err)
				}
			}
		}
	}
	for _, hook := range snap.afterHooks {
		hook(entry)
	}
}

// writeFormatted escreve bytes já formatados diretamente no writer do transporte.
func writeFormatted(t Transport, data []byte) error {
	switch tr := t.(type) {
	case *WriterTransport:
		_, err := tr.Writer.Write(data)
		return err
	case *ConsoleTransport:
		out := os.Stdout
		if tr.ToStdErr {
			out = os.Stderr
		}
		_, err := out.Write(data)
		return err
	case *FileTransport:
		_, err := tr.File.Write(data)
		return err
	case *LumberjackTransport:
		_, err := tr.Logger.Write(data)
		return err
	default:
		// Fallback: usa WriteLog normal (ignora formatter customizado)
		return t.WriteLog(&Entry{Message: string(data)})
	}
}

// log envia a entry para todos os transportes cujo nível mínimo seja compatível.
func (l *Logger) log(level Level, message string) {
	snap := l.snapshot()
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
	}
	if snap.stacktrace.Enabled && snap.stacktrace.Levels[level] {
		if entry.Fields == nil {
			entry.Fields = make(map[string]interface{})
		}
		entry.Fields["stacktrace"] = string(debug.Stack())
	}
	dispatchEntry(snap, &entry, nil)
}

// Debug registra uma mensagem no nível DEBUG.
func (l *Logger) Debug(message string) {
	l.log(DEBUG, message)
}

// Info registra uma mensagem no nível INFO.
func (l *Logger) Info(message string) {
	l.log(INFO, message)
}

// Warn registra uma mensagem no nível WARN.
func (l *Logger) Warn(message string) {
	l.log(WARN, message)
}

// Error registra uma mensagem no nível ERROR.
func (l *Logger) Error(message string) {
	l.log(ERROR, message)
}

// Fatal registra uma mensagem no nível ERROR, inclui stacktrace e encerra a aplicação.
func (l *Logger) Fatal(message string, fields ...map[string]any) {
	var flds map[string]any
	if len(fields) > 0 {
		flds = fields[0]
	}
	if flds == nil {
		flds = make(map[string]any)
	}
	flds["stacktrace"] = string(debug.Stack())
	l.logWithFields(ERROR, message, flds)
	os.Exit(1)
}

// Panic registra uma mensagem no nível ERROR, inclui stacktrace e faz panic.
func (l *Logger) Panic(message string, fields ...map[string]any) {
	var flds map[string]any
	if len(fields) > 0 {
		flds = fields[0]
	}
	if flds == nil {
		flds = make(map[string]any)
	}
	flds["stacktrace"] = string(debug.Stack())
	l.logWithFields(ERROR, message, flds)
	panic(message)
}

// logWithFields é usada internamente por EntryBuilder.
func (l *Logger) logWithFields(level Level, message string, fields map[string]interface{}) {
	snap := l.snapshot()
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}
	if snap.stacktrace.Enabled && snap.stacktrace.Levels[level] {
		if entry.Fields == nil {
			entry.Fields = make(map[string]interface{})
		}
		entry.Fields["stacktrace"] = string(debug.Stack())
	}
	dispatchEntry(snap, &entry, nil)
}

// ComFields permite adicionar metadata/contexto extra ao log.
func (l *Logger) ComFields(fields map[string]interface{}) *EntryBuilder {
	return &EntryBuilder{logger: l, fields: fields}
}

// EntryBuilder permite construir logs com metadata extra e formatter customizado.
type EntryBuilder struct {
	logger    *Logger
	fields    map[string]interface{}
	formatter Formatter
}

// WithFormatter permite sobrescrever o formatter para este log.
func (l *Logger) WithFormatter(formatter Formatter) *EntryBuilder {
	return &EntryBuilder{logger: l, formatter: formatter}
}

func (b *EntryBuilder) Debug(msg string) {
	b.logger.logWithFieldsCustomFormatter(DEBUG, msg, b.fields, b.formatter)
}
func (b *EntryBuilder) Info(msg string) {
	b.logger.logWithFieldsCustomFormatter(INFO, msg, b.fields, b.formatter)
}
func (b *EntryBuilder) Warn(msg string) {
	b.logger.logWithFieldsCustomFormatter(WARN, msg, b.fields, b.formatter)
}
func (b *EntryBuilder) Error(msg string) {
	b.logger.logWithFieldsCustomFormatter(ERROR, msg, b.fields, b.formatter)
}

// logWithFieldsCustomFormatter permite sobrescrever o formatter por mensagem (thread-safe).
func (l *Logger) logWithFieldsCustomFormatter(level Level, message string, fields map[string]interface{}, formatter Formatter) {
	snap := l.snapshot()
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}
	dispatchEntry(snap, &entry, formatter)
}

// LoggerConfig permite inicializar o logger de forma dinâmica.
type LoggerConfig struct {
	Transports []TransportConfig
}

type TransportConfig struct {
	Type      string         // "console", "file", etc
	Level     string         // "INFO", "DEBUG", ...
	Formatter string         // "text", "json"
	Options   map[string]any // opções específicas (ex: path para arquivo)
}

// NewLoggerFromConfig cria um Logger a partir de uma configuração dinâmica.
func NewLoggerFromConfig(cfg LoggerConfig) (*Logger, error) {
	logger := &Logger{}
	for _, tcfg := range cfg.Transports {
		var formatter Formatter
		switch tcfg.Formatter {
		case "json":
			formatter = &JSONFormatter{}
		default:
			formatter = &TextFormatter{}
		}
		level := ParseLevel(tcfg.Level)
		switch tcfg.Type {
		case "console":
			toStdErr := false
			if v, ok := tcfg.Options["stderr"].(bool); ok {
				toStdErr = v
			}
			logger.AddTransport(&ConsoleTransport{
				Level:     level,
				Formatter: formatter,
				ToStdErr:  toStdErr,
			})
		case "file":
			path, _ := tcfg.Options["path"].(string)
			ft, err := NewFileTransport(path, level, formatter)
			if err != nil {
				return nil, err
			}
			logger.AddTransport(ft)
		default:
			return nil, fmt.Errorf("lazylog: unknown transport type %q", tcfg.Type)
		}
	}
	return logger, nil
}

// LoadLoggerConfigJSON carrega configuração do logger de um arquivo JSON.
func LoadLoggerConfigJSON(path string) (LoggerConfig, error) {
	var cfg LoggerConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

// LoadLoggerConfigYAML carrega configuração do logger de um arquivo YAML.
func LoadLoggerConfigYAML(path string) (LoggerConfig, error) {
	var cfg LoggerConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

// logWithContext permite logar com context.Context, extraindo informações relevantes.
func (l *Logger) logWithContext(ctx context.Context, level Level, message string, fields map[string]interface{}) {
	snap := l.snapshot()
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}
	// Suporte a context key customizada e string
	for _, key := range []any{CtxKey("trace_id"), "trace_id"} {
		if v := ctx.Value(key); v != nil {
			if entry.Fields == nil {
				entry.Fields = make(map[string]interface{})
			}
			entry.Fields["trace_id"] = v
			break
		}
	}
	dispatchEntry(snap, &entry, nil)
}

// API pública para logar com contexto
func (l *Logger) InfoCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logWithContext(ctx, INFO, msg, fields)
}
func (l *Logger) DebugCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logWithContext(ctx, DEBUG, msg, fields)
}
func (l *Logger) WarnCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logWithContext(ctx, WARN, msg, fields)
}
func (l *Logger) ErrorCtx(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logWithContext(ctx, ERROR, msg, fields)
}

// CtxKey é o tipo exportado para chaves de contexto do lazylog.
type CtxKey string

// ChildLogger permite criar um logger derivado com campos fixos (contexto).
type ChildLogger struct {
	parent *Logger
	fields map[string]any
}

// WithFields retorna um logger derivado com contexto fixo.
func (l *Logger) WithFields(fields map[string]any) *ChildLogger {
	return &ChildLogger{parent: l, fields: fields}
}

func (c *ChildLogger) Debug(msg string, fields ...map[string]any) {
	c.logWithMergedFields(DEBUG, msg, fields...)
}
func (c *ChildLogger) Info(msg string, fields ...map[string]any) {
	c.logWithMergedFields(INFO, msg, fields...)
}
func (c *ChildLogger) Warn(msg string, fields ...map[string]any) {
	c.logWithMergedFields(WARN, msg, fields...)
}
func (c *ChildLogger) Error(msg string, fields ...map[string]any) {
	c.logWithMergedFields(ERROR, msg, fields...)
}
func (c *ChildLogger) logWithMergedFields(level Level, msg string, fields ...map[string]any) {
	merged := make(map[string]any)
	for k, v := range c.fields {
		merged[k] = v
	}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			merged[k] = v
		}
	}
	c.parent.logWithFields(level, msg, merged)
}
