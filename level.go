package logger

import "fmt"

type Level int

const (
	LevelError Level = iota
	LevelInfo
	LevelDebug
	//
	LevelDefault //use parent's level
)

func (l Level) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelDefault:
		return "DEFAULT"
	default:
	}
	return fmt.Sprintf("LEVEL(%v)", int(l))
}
