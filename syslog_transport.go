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
	bytes, err := s.Formatter.Format(entry)
	if err != nil {
		return err
	}
	return s.Writer.Info(string(bytes))
}

func (s *SyslogTransport) MinLevel() Level {
	return s.Level
}
