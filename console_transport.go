package lazylog

import (
	"os"
)

// ConsoleTransport escreve logs no stdout ou stderr.
type ConsoleTransport struct {
	Level     Level
	Formatter Formatter
	ToStdErr  bool // Se true, escreve no stderr; senão, no stdout
}

func (c *ConsoleTransport) WriteLog(entry *Entry) error {
	out := os.Stdout
	if c.ToStdErr {
		out = os.Stderr
	}
	formatter := c.Formatter
	if formatter == nil {
		formatter = &TextFormatter{}
	}
	bytes, err := formatter.Format(entry)
	if err != nil {
		_, err2 := out.Write([]byte(entry.Timestamp.Format("2006-01-02T15:04:05Z07:00") + " [" + entry.Level.String() + "] " + entry.Message + "\n"))
		return err2
	}
	_, err = out.Write(bytes)
	return err
}

func (c *ConsoleTransport) MinLevel() Level {
	return c.Level
}
