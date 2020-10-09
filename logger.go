package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelError
)

type ILogger interface {
	NewLogger(name string) ILogger //name is valid basename used to construct the path (any simple string/int value)
	Parent() ILogger
	Path() string //path is [/<parent path>/<name>]
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	WithStream(ILogStream) ILogger
}

func NewLogger(name string) ILogger {
	return &logger{
		parent:  nil,
		name:    name,
		path:    "/" + name,
		streams: []ILogStream{},
	}
}

func NewTerminalLogger(name string) ILogger {
	return &logger{
		parent: nil,
		name:   name,
		path:   "/" + name,
		streams: []ILogStream{
			logStream{
				level:   LogLevelDebug,
				encoder: logEncoderTerminal{},
				writer:  logWriterFile{file: os.Stderr},
			},
		},
	}
}

type logger struct {
	parent  ILogger
	name    string
	path    string
	streams []ILogStream
}

func (l *logger) NewLogger(name string) ILogger {
	child := &logger{
		parent:  l,
		name:    name,
		path:    l.path + "/" + name,
		streams: []ILogStream{},
	}
	return child
}

func (l *logger) Name() string    { return l.name }
func (l *logger) Parent() ILogger { return l.parent }
func (l *logger) Path() string    { return l.path }

func (l *logger) WithStream(s ILogStream) ILogger {
	if s != nil {
		l.streams = append(l.streams, s)
	}
	return l
}

func (l *logger) log(caller Caller, level LogLevel, msg string) {
	//send to all streams
	for _, s := range l.streams {
		s.Log(LogRecord{
			Caller:    caller,
			Timestamp: time.Now(),
			Logger:    l,
			Level:     level,
			Message:   msg,
		})
	}
}

func (l *logger) Log(level LogLevel, msg string) {
	l.log(GetCaller(1), level, msg)
}

func (l logger) Logf(level LogLevel, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), level, msg)
}

func (l logger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), LogLevelDebug, msg)
}

func (l logger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), LogLevelInfo, msg)
}

func (l logger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), LogLevelError, msg)
}

type ILogStream interface {
	WithWriter(ILogWriter) ILogStream
	WithEncoder(ILogEncoder) ILogStream
	Log(record LogRecord)
}

func NewLogStream() ILogStream {
	return logStream{}
}

type logStream struct {
	level   LogLevel
	encoder ILogEncoder
	writer  ILogWriter
}

func (s logStream) WithWriter(w ILogWriter) ILogStream {
	s.writer = w
	return s
}

func (s logStream) WithEncoder(e ILogEncoder) ILogStream {
	s.encoder = e
	return s
}

func (s logStream) Log(record LogRecord) {
	if s.encoder != nil && s.writer != nil && record.Level >= s.level {
		s.writer.Write(s.encoder.Encode(record))
	}
}

type LogRecord struct {
	Timestamp time.Time
	Logger    ILogger
	Caller    Caller
	Level     LogLevel
	Message   string
	Data      map[string]interface{}
}

type ILogEncoder interface {
	Encode(r LogRecord) []byte
}

type logEncoderTerminal struct{}

func (e logEncoderTerminal) Encode(r LogRecord) []byte {
	return []byte(fmt.Sprintf("%s %5.5s %s: %s",
		r.Timestamp.Format("2006-01-02 15:04:05.000"),
		"DEBUG",
		r.Caller,
		r.Message,
	))
}

type ILogWriter interface {
	Write([]byte)
}

type logWriterFile struct {
	file *os.File
}

func (w logWriterFile) Write(record []byte) {
	//replace any newlines with ';' then append a newline only at the end
	w.file.Write([]byte(strings.ReplaceAll(string(record), "\n", ";") + "\n"))
}
