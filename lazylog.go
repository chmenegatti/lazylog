package lazylog

import (
	"time"
)

type Logger struct {
	transports []Transport
}

// NewLogger cria um logger com zero ou mais transportes.
func NewLogger(transports ...Transport) *Logger {
	return &Logger{
		transports: transports,
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

// log envia a entry para todos os transportes cujo nível mínimo seja compatível.
func (l *Logger) log(level Level, message string) {
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
	}
	for _, t := range l.transports {
		if level >= t.MinLevel() {
			t.WriteLog(&entry)
		}
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

func (l *Logger) logWithFields(level Level, message string, fields map[string]interface{}) {
	entry := Entry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}
	for _, t := range l.transports {
		if level >= t.MinLevel() {
			t.WriteLog(&entry)
		}
	}
}

// ComFields permite adicionar metadata/contexto extra ao log.
func (l *Logger) ComFields(fields map[string]interface{}) *EntryBuilder {
	return &EntryBuilder{logger: l, fields: fields}
}

// EntryBuilder permite construir logs com metadata extra.
type EntryBuilder struct {
	logger *Logger
	fields map[string]interface{}
}

func (b *EntryBuilder) Debug(msg string) { b.logger.logWithFields(DEBUG, msg, b.fields) }
func (b *EntryBuilder) Info(msg string)  { b.logger.logWithFields(INFO, msg, b.fields) }
func (b *EntryBuilder) Warn(msg string)  { b.logger.logWithFields(WARN, msg, b.fields) }
func (b *EntryBuilder) Error(msg string) { b.logger.logWithFields(ERROR, msg, b.fields) }

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
