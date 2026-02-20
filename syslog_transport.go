package lazylog

import (
	"log/syslog"
)

// SyslogTransport envia logs para o syslog do sistema.
type SyslogTransport struct {
	Writer    *syslog.Writer
	Level     Level
	Formatter Formatter
}

func NewSyslogTransport(priority syslog.Priority, tag string, level Level, formatter Formatter) (*SyslogTransport, error) {
	writer, err := syslog.New(priority, tag)
	if err != nil {
		return nil, err
	}
	return &SyslogTransport{
		Writer:    writer,
		Level:     level,
		Formatter: formatter,
	}, nil
}

func (s *SyslogTransport) WriteLog(entry *Entry) error {
	formatter := s.Formatter
	if formatter == nil {
		formatter = &TextFormatter{}
	}
	bytes, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	msg := string(bytes)
	switch entry.Level {
	case DEBUG:
		return s.Writer.Debug(msg)
	case WARN:
		return s.Writer.Warning(msg)
	case ERROR:
		return s.Writer.Err(msg)
	default:
		return s.Writer.Info(msg)
	}
}

func (s *SyslogTransport) MinLevel() Level {
	return s.Level
}

// Close fecha o writer do syslog.
func (s *SyslogTransport) Close() error {
	return s.Writer.Close()
}
