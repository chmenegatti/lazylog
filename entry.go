package lazylog

import "time"

type Entry struct {
	Level     Level
	Timestamp time.Time
	Message   string
	Fields    map[string]interface{} // Para metadata/contexto extra
}
