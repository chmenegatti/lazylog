package lazylog

import (
	"strings"
	"sync"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	levelMu    sync.RWMutex
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
	levelMu.Lock()
	defer levelMu.Unlock()
	levelNames[value] = name
	levelValues[name] = value
}

func (l Level) String() string {
	levelMu.RLock()
	defer levelMu.RUnlock()
	if name, ok := levelNames[l]; ok {
		return name
	}
	return "UNKNOWN"
}

func ParseLevel(lvl string) Level {
	levelMu.RLock()
	defer levelMu.RUnlock()
	if v, ok := levelValues[strings.ToUpper(lvl)]; ok {
		return v
	}
	return INFO // Default to INFO if the level is unknown
}
