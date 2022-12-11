package logger

import (
	"fmt"
	"strings"
	"time"
)

// Logger does logging
// it can only be created from a INamedLogger which gives the structure to control log levels by name during run-time
//
//	but that actually only control the writer's level, which is used by all loggers with the same name
//
// after creating a logger, the logger level can be changed in code, e.g. by a trigger condition
//
//	but Logger level cannot change from outside. By default, Logger uses the level of the named logger from where
//	it was created
type Logger interface {
	New(name string) Logger

	//SetXxx affects the named logger and all sub named loggers
	//but does not change loggers already created in those names, just the names when they are used to create new loggers
	SetWriter(newWriter IWriter) //sets the writer on this and all child loggers
	SetLevel(newLevel Level)     //sets the level on named logger, for this logger and all NEW child loggers and existing child loggers that did not use WithLevel()

	//WithXxx creates a copy of the logger with the new settings...
	WithLevel(Level) Logger //only affects this new logger (can use LevelDefault to reset to named level and allow external control)
	With(name string, value interface{}) Logger

	Name() string
	Names() []string
	Level() Level

	Log(level Level, msg string)
	Error(msg string)
	Info(msg string)
	Debug(msg string)

	Logf(level Level, format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type logger struct {
	named *named
	level Level
	data  map[string]interface{}
}

func (l logger) New(name string) Logger {
	return l.named.New(name)
}

func (l logger) SetWriter(newWriter IWriter) {
	l.named.setWriter(newWriter)
}

func (l logger) SetLevel(newLevel Level) {
	l.level = LevelDefault //clear own setting and use named's level...
	l.named.setLevel(newLevel)
}

func (l logger) Name() string    { return l.named.name }
func (l logger) Names() []string { return l.named.names }

func (l logger) Level() Level {
	if l.level < LevelDefault {
		return l.level //fall through to use named logger's level
	}
	return l.named.level
}

func (l logger) String() string { return l.named.name }

func (l logger) WithLevel(newLevel Level) Logger {
	l.level = newLevel
	return l
}

func (l logger) With(name string, value interface{}) Logger {
	d := l.data
	l.data = map[string]interface{}{}
	for n, v := range d {
		l.data[n] = v
	}
	l.data[name] = value
	return l
}

func (l logger) Log(level Level, msg string) { l.log(3, level, msg) }
func (l logger) Error(msg string)            { l.log(3, LevelError, msg) }
func (l logger) Info(msg string)             { l.log(3, LevelInfo, msg) }
func (l logger) Debug(msg string)            { l.log(3, LevelDebug, msg) }

func (l logger) Logf(level Level, format string, args ...interface{}) {
	l.logf(4, level, format, args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.logf(4, LevelError, format, args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.logf(4, LevelInfo, format, args...)
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.logf(4, LevelDebug, format, args...)
}

func (l logger) log(depth int, level Level, msg string) {
	if level <= l.Level() {
		l.named.writer.Write(
			Record{
				Caller:    GetCaller(depth),
				Timestamp: time.Now(),
				Logger:    l,
				Level:     level,
				Message:   strings.ReplaceAll(msg, "\n", "; "),
				Data:      l.data,
			},
		)
	}
}

func (l logger) logf(depth int, level Level, format string, args ...interface{}) {
	if level <= l.Level() {
		l.log(depth, level, fmt.Sprintf(format, args...))
	}
}
