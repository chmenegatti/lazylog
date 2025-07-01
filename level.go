package lazylog

import "strings"

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	levelNames = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}
	levelValues = map[string]Level{
		"DEBUG": DEBUG,
		"INFO":  INFO,
		"WARN":  WARN,
		"ERROR": ERROR,
	}
)

// RegisterLevel permite registrar um novo nível de log customizado.
func RegisterLevel(name string, value Level) {
	levelNames[value] = name
	levelValues[name] = value
}

func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return "UNKNOWN"
}

func ParseLevel(lvl string) Level {
	if v, ok := levelValues[strings.ToUpper(lvl)]; ok {
		return v
	}
	return INFO // Default to INFO if the level is unknown
}
