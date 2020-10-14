package logger

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelError
)

type ILogger interface {
	NewLogger(name string) ILogger //name must be valid pathname
	Parent() ILogger
	Children() map[string]ILogger

	Path() string //path is [/<parent path>/<name>]
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	WithStream(ILogStream) ILogger
}

func Top() ILogger {
	return top
}

func ForThisPackage() ILogger {
	c := GetCaller(1)
	return top.NewLogger(c.Package())
}

//Terminal is a stream for logging to stderr
func Terminal(level LogLevel) ILogStream {
	return logStream{
		level:   level,
		encoder: logEncoderTerminal{},
		writer:  logWriterFile{file: os.Stderr},
	}
}

// func NewTerminalLogger(name string) ILogger {
// 	return &logger{
// 		parent: nil,
// 		name:   name,
// 		path:   "/" + name,
// 		streams: []ILogStream{
// 			logStream{
// 				level:   LogLevelDebug,
// 				encoder: logEncoderTerminal{},
// 				writer:  logWriterFile{file: os.Stderr},
// 			},
// 		},
// 	}
// }

type logger struct {
	sync.Mutex
	parent   ILogger
	children map[string]ILogger
	name     string
	path     string
	streams  []ILogStream
}

func (l *logger) NewLogger(name string) ILogger {
	l.Lock()
	defer l.Unlock()

	names := strings.SplitN(path.Clean(name), "/", 2)
	for len(names) > 1 && len(names[0]) == 0 {
		names = strings.SplitN(names[1], "/", 2) //skip empty name
	}
	if len(names) == 0 || len(names[0]) == 0 {
		return l
	}

	//get/create child
	child, ok := l.children[names[0]]
	if !ok {
		path := names[0]
		if l.path != "" {
			path = l.path + "/" + names[0]
		}

		child = &logger{
			parent:   l,
			children: map[string]ILogger{},
			name:     names[0],
			path:     path,
			streams:  nil, //[]ILogStream{},
		}
		l.children[names[0]] = child
		fmt.Printf("created logger(%s)\n", child.Path())
	}

	if len(names) > 1 {
		return child.NewLogger(names[1])
	}
	return child
}

func (l *logger) Name() string                 { return l.name }
func (l *logger) Parent() ILogger              { return l.parent }
func (l *logger) Children() map[string]ILogger { return l.children }
func (l *logger) Path() string                 { return l.path }

func (l *logger) WithStream(s ILogStream) ILogger {
	if l != nil {
		if s != nil {
			if l.streams == nil {
				l.streams = []ILogStream{s}
			} else {
				l.streams = append(l.streams, s)
			}
		}
	}
	return l
}

func (l *logger) log(caller Caller, level LogLevel, msg string) {
	//send to all streams
	//if no streams in this logger, use parent streams
	if l.streams == nil && l.parent != nil {
		if p, ok := l.parent.(*logger); ok {
			p.log(caller, level, msg)
		}
	} else {
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
}

func (l *logger) Log(level LogLevel, msg string) {
	l.log(GetCaller(1), level, msg)
}

func (l *logger) Logf(level LogLevel, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), level, msg)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), LogLevelDebug, msg)
}

func (l *logger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), LogLevelInfo, msg)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(GetCaller(2), LogLevelError, msg)
}

func (l *logger) debugOn() {
	//let all streams write to debug

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
	return []byte(fmt.Sprintf("%s %5.5s %25.5s: %s",
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

var (
	top *logger
)

func init() {
	//top by default has no streams and logs nothing
	//DebugOn() can be used to let every logger write to the terminal
	//DebugPackage(<name>) can be used to let a package write to the terminal
	top = &logger{
		parent:   nil,
		children: map[string]ILogger{},
		name:     "",
		path:     "",
		streams:  nil, //[]ILogStream{},
	}
}
