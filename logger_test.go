package logger_test

import (
	"testing"

	"github.com/go-msvc/logger"
)

type logEncoder struct {
	t *testing.T
}

func (e logEncoder) Encode(r logger.LogRecord) []byte {
	e.t.Logf("encoding %+v", r)
	return []byte(r.Message)
}

type logWriter struct {
	t    *testing.T
	logs string
}

func (w *logWriter) Write(r []byte) {
	w.t.Logf("writing %v", r)
	w.logs += string(r)
}

func TestLogger(t *testing.T) {
	w := &logWriter{t: t}
	e := logEncoder{t: t}
	l := logger.Top().NewLogger("test").
		WithStream(logger.NewLogStream().WithWriter(w).WithEncoder(e))

	w.logs = ""
	l.Debugf("debug")
	if w.logs != "debug" {
		t.Fatalf("logged: \"%s\" != \"debug\"", w.logs)
	}
}
