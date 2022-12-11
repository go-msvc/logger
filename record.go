package logger

import "time"

type Record struct {
	Timestamp time.Time
	Logger    Logger
	Caller    Caller
	Level     Level
	Message   string
	Data      map[string]interface{}
}
