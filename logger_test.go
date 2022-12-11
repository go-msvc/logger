package logger_test

import (
	"strings"
	"testing"

	"github.com/go-msvc/logger"
	"github.com/go-msvc/logger/fakelib"
)

type testWriter struct {
	records []logger.Record
}

func (w *testWriter) Write(r logger.Record) {
	w.records = append(w.records, r)
}

func (w *testWriter) Reset() {
	w.records = []logger.Record{}
}

func (w *testWriter) assert(t *testing.T, recordIndex int, funcName, loggerName, level, msg string, data map[string]interface{}) {
	t.Logf("assert(%d,%s,%s,%s,%s,%+v)", recordIndex, funcName, loggerName, level, msg, data)
	if recordIndex >= len(w.records) {
		t.Fatalf("record[%d] not logged", recordIndex)
	}
	r := w.records[recordIndex]
	if r.Caller.Function() != funcName {
		t.Fatalf("[%d] func=%s != %s", recordIndex, r.Caller.Function(), funcName)
	}
	if r.Logger.Name() != loggerName {
		t.Fatalf("[%d] loggerName=%s != %s", recordIndex, r.Logger.Name(), loggerName)
	}
	if r.Level.String() != level {
		t.Fatalf("[%d] level=%s != %s", recordIndex, r.Level, level)
	}
	if r.Message != msg {
		t.Fatalf("[%d] msg=%s != %s", recordIndex, r.Message, msg)
	}
	for n, v := range data {
		if r.Data[n] != v {
			t.Fatalf("[%d] data: %s=(%T)%v != (%T)%v", recordIndex, n, r.Data[n], r.Data[n], v, v)
		}
	}
}

func TestOne(t *testing.T) {
	//by default, everything is on ERROR so switch to debug for own logs
	w := &testWriter{records: []logger.Record{}}
	l := logger.New().WithLevel(logger.LevelDebug)
	l.SetWriter(w)

	l.Debugf("123")
	l = l.With("email", "a@b.c")
	l.Debugf("456")

	l = logger.New().WithLevel(logger.LevelDebug)
	l.Debugf("789")
	l = l.With("email", "j@k.l")
	l.Debugf("101")

	w.assert(t, 0, "TestOne", "github.com/go-msvc/logger_test", "DEBUG", "123", map[string]interface{}{})
	w.assert(t, 1, "TestOne", "github.com/go-msvc/logger_test", "DEBUG", "456", map[string]interface{}{"email": "a@b.c"})
	w.assert(t, 2, "TestOne", "github.com/go-msvc/logger_test", "DEBUG", "789", map[string]interface{}{})
	w.assert(t, 3, "TestOne", "github.com/go-msvc/logger_test", "DEBUG", "101", map[string]interface{}{"email": "j@k.l"})
}

func TestFakelibDefault(t *testing.T) {
	//a program may traverse the tree of named loggers
	t.Logf("ALL LOGGERS:")
	for _, l := range logger.All() {
		showLogger(t, l)
	}

	//use testWriter to collect all logs for evaluation
	w := &testWriter{records: []logger.Record{}}
	logger.SetGlobalWriter(w)

	//a main program by default should be able to import fakelib and only get ERROR logs from it
	fakelib.Fake()
	for _, r := range w.records {
		// t.Logf("record[%d]: %+v", i, r)
		if r.Level > logger.LevelError {
			t.Fatalf("got non-error log: %+v", r)
		}
	}
	w.Reset()

	//a main program can switch on logger for a library when required
	logger.Named("github.com/go-msvc/logger/fakelib").SetLevel(logger.LevelInfo)
	t.Logf("ALL LOGGERS:")
	for _, l := range logger.All() {
		showLogger(t, l)
	}

	//now running fakelib func will log at level info
	fakelib.Fake()
	infoCount := 0
	for _, r := range w.records {
		// t.Logf("record[%d]: %+v", i, r)
		if r.Level > logger.LevelInfo {
			t.Fatalf("got >info log: %+v", r)
		}
		if r.Level == logger.LevelInfo {
			infoCount++
		}
	}
	if infoCount != 5 {
		t.Fatalf("infoCount=%d", infoCount)
	}
	w.Reset()

	//or switch to debug:
	logger.Named("github.com/go-msvc/logger/fakelib").SetLevel(logger.LevelDebug)
	t.Logf("ALL LOGGERS:")
	for _, l := range logger.All() {
		showLogger(t, l)
	}

	//now running fakelib func will log at level info
	fakelib.Fake()
	debugCount := 0
	for _, r := range w.records {
		// t.Logf("record[%d]: %+v", i, r)
		if r.Level == logger.LevelInfo {
			debugCount++
		}
	}
	if debugCount != 5 {
		t.Fatalf("debugCount=%d", debugCount)
	}
	w.Reset()

}

func showLogger(t *testing.T, l logger.INamed) {
	t.Logf("%s: %s", l.Level(), strings.Join(l.Names(), "/"))
	for _, l := range l.All() {
		showLogger(t, l)
	}
}
