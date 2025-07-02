package lazylog

import (
	"context"
	"runtime/debug"
	"time"
)

type Hook func(entry *Entry)

// Hook de erro de transporte
// TransportErrorHook é chamado quando um transporte falha ao gravar.
type TransportErrorHook func(entry *Entry, transport Transport, err error)

// StacktraceConfig permite ativar stacktrace automático para níveis específicos.
type StacktraceConfig struct {
	Enabled bool
	Levels  map[Level]bool // Níveis que devem incluir stacktrace
}

// Logger agora pode ter configuração de stacktrace.
type Logger struct {
	transports  []Transport
	BeforeHooks []Hook
	AfterHooks  []Hook
	ErrorHooks  []TransportErrorHook
	Stacktrace  StacktraceConfig
}

// NewLogger cria um logger com zero ou mais transportes.
func NewLogger(transports ...Transport) *Logger {
	return &Logger{
		transports: transports,
	}
}

// EnableStacktrace ativa stacktrace automático para os níveis informados.
func (l *Logger) EnableStacktrace(levels ...Level) {
	l.Stacktrace.Enabled = true
	l.Stacktrace.Levels = make(map[Level]bool)
	for _, lvl := range levels {
		l.Stacktrace.Levels[lvl] = true
	}
}

// AddTransport adiciona um novo transporte ao logger.
func (l *Logger) AddTransport(t Transport) {
	l.transports = append(l.transports, t)
}

// RemoveTransport remove um transporte do logger (por comparação de ponteiro).
func (l *Logger) RemoveTransport(t Transport) {
	for i, tr := range l.transports {
		if tr == t {
			l.transports = append(l.transports[:i], l.transports[i+1:]...)
			return
		}
	}
}

// AddHook adiciona um hook para ser executado antes ou depois do log.
func (l *Logger) AddHook(hook Hook, before bool) {
	if before {
		l.BeforeHooks = append(l.BeforeHooks, hook)
	} else {
		l.AfterHooks = append(l.AfterHooks, hook)
	}
}

// AddErrorHook adiciona um hook para erros de transporte.
func (l *Logger) AddErrorHook(hook TransportErrorHook) {
	l.ErrorHooks = append(l.ErrorHooks, hook)
}

// log envia a entry para todos os transportes cujo nível mínimo seja compatível.
func (l *Logger) log(level Level, message string) {
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
	}
	if l.Stacktrace.Enabled && l.Stacktrace.Levels[level] {
		if entry.Fields == nil {
			entry.Fields = make(map[string]interface{})
		}
		entry.Fields["stacktrace"] = string(debug.Stack())
	}
	for _, hook := range l.BeforeHooks {
		hook(&entry)
	}
	for _, t := range l.transports {
		if level >= t.MinLevel() {
			err := t.WriteLog(&entry)
			if err != nil {
				for _, eh := range l.ErrorHooks {
					eh(&entry, t, err)
				}
			}
		}
	}
	for _, hook := range l.AfterHooks {
		hook(&entry)
	}
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

// logWithFields é usada internamente por EntryBuilder e pode ser exportada se desejado.
// Para evitar o aviso de função não utilizada, pode-se adicionar um comentário '//lint:ignore U1000 used by builder' ou exportar se for útil externamente.
//
//lint:ignore U1000 used by EntryBuilder
func (l *Logger) logWithFields(level Level, message string, fields map[string]interface{}) {
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}
	if l.Stacktrace.Enabled && l.Stacktrace.Levels[level] {
		if entry.Fields == nil {
			entry.Fields = make(map[string]interface{})
		}
		entry.Fields["stacktrace"] = string(debug.Stack())
	}
	for _, hook := range l.BeforeHooks {
		hook(&entry)
	}
	for _, t := range l.transports {
		if level >= t.MinLevel() {
			err := t.WriteLog(&entry)
			if err != nil {
				for _, eh := range l.ErrorHooks {
					eh(&entry, t, err)
				}
			}
		}
	}
	for _, hook := range l.AfterHooks {
		hook(&entry)
	}
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

// logWithFieldsCustomFormatter permite sobrescrever o formatter por mensagem.
func (l *Logger) logWithFieldsCustomFormatter(level Level, message string, fields map[string]interface{}, formatter Formatter) {
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}
	for _, hook := range l.BeforeHooks {
		hook(&entry)
	}
	for _, t := range l.transports {
		if level >= t.MinLevel() {
			if formatter != nil {
				// Se o transporte for WriterTransport, ConsoleTransport, FileTransport ou LumberjackTransport, sobrescreve o formatter temporariamente
				switch tr := t.(type) {
				case *WriterTransport:
					orig := tr.Formatter
					tr.Formatter = formatter
					_ = tr.WriteLog(&entry)
					tr.Formatter = orig
				case *ConsoleTransport:
					orig := tr.Formatter
					tr.Formatter = formatter
					_ = tr.WriteLog(&entry)
					tr.Formatter = orig
				case *FileTransport:
					orig := tr.Formatter
					tr.Formatter = formatter
					_ = tr.WriteLog(&entry)
					tr.Formatter = orig
				case *LumberjackTransport:
					orig := tr.Formatter
					tr.Formatter = formatter
					_ = tr.WriteLog(&entry)
					tr.Formatter = orig
				default:
					_ = t.WriteLog(&entry)
				}
			} else {
				_ = t.WriteLog(&entry)
			}
		}
	}
	for _, hook := range l.AfterHooks {
		hook(&entry)
	}
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
		}
	}
	return logger, nil
}

// logWithContext permite logar com context.Context, extraindo informações relevantes.
func (l *Logger) logWithContext(ctx context.Context, level Level, message string, fields map[string]interface{}) {
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}
	// Suporte a context key customizada e string
	for _, key := range []any{ctxKey("trace_id"), "trace_id"} {
		if v := ctx.Value(key); v != nil {
			if entry.Fields == nil {
				entry.Fields = make(map[string]interface{})
			}
			entry.Fields["trace_id"] = v
			break
		}
	}
	for _, hook := range l.BeforeHooks {
		hook(&entry)
	}
	for _, t := range l.transports {
		if level >= t.MinLevel() {
			err := t.WriteLog(&entry)
			if err != nil {
				for _, eh := range l.ErrorHooks {
					eh(&entry, t, err)
				}
			}
		}
	}
	for _, hook := range l.AfterHooks {
		hook(&entry)
	}
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

type ctxKey string
