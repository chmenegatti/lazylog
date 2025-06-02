package lazylog

import "time"

type Entry struct {
	Level     Level
	Timestamp time.Time
	Message   string
}
