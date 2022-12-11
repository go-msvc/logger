package logger_test

import (
	"testing"

	"github.com/go-msvc/logger"
)

func test(t *testing.T, w *testWriter, l logger.Logger) {
	l.Debugf("123") //line 10
	l.Infof("123")  //line 11
	l.Errorf("123") //line 12
	l = l.With("email", "a@b.c")
	l.Debugf("456") //line 14

	l = l.New(l.Name()).WithLevel(logger.LevelDebug)
	l.Debugf("789") //line 17
	l = l.With("email", "j@k.l")
	l.Debugf("101") //line 19
}

func TestCodeWriterGlobal(t *testing.T) {
	w := &testWriter{records: []logger.Record{}}
	cw := logger.NewCodeWriter(w, logger.LevelDebug)
	l := logger.New().WithLevel(logger.LevelDebug)
	l.SetWriter(cw)
	test(t, w, l)

	funcName := "test" //where it was writter
	w.assert(t, 0, funcName, l.Name(), "DEBUG", "123", map[string]interface{}{})
	w.assert(t, 1, funcName, l.Name(), "INFO", "123", map[string]interface{}{})
	w.assert(t, 2, funcName, l.Name(), "ERROR", "123", map[string]interface{}{})
	w.assert(t, 3, funcName, l.Name(), "DEBUG", "456", map[string]interface{}{"email": "a@b.c"})
	w.assert(t, 4, funcName, l.Name(), "DEBUG", "789", map[string]interface{}{})
	w.assert(t, 5, funcName, l.Name(), "DEBUG", "101", map[string]interface{}{"email": "j@k.l"})
}

func TestCodeWriterFile(t *testing.T) {
	w := &testWriter{records: []logger.Record{}}
	cw := logger.NewCodeWriter(w, logger.LevelDebug)
	cw.SetFileLevel("github.com/go-msvc/logger_test/code-writer_test.go", logger.LevelInfo)
	l := logger.New().WithLevel(logger.LevelDebug)
	logger.SetGlobalWriter(cw)
	//this does nothing: logger.SetGlobalLevel(logger.LevelError)
	//because we control levels in the writer!
	test(t, w, l)

	funcName := "test" //where it was writter
	//w.assert(t, 0, funcName, l.Name(), "DEBUG", "123", map[string]interface{}{})
	w.assert(t, 0, funcName, l.Name(), "INFO", "123", map[string]interface{}{})
	w.assert(t, 1, funcName, l.Name(), "ERROR", "123", map[string]interface{}{})
	//w.assert(t, 3, funcName, l.Name(), "DEBUG", "456", map[string]interface{}{"email": "a@b.c"})
	//w.assert(t, 4, funcName, l.Name(), "DEBUG", "789", map[string]interface{}{})
	//w.assert(t, 2, funcName, l.Name(), "DEBUG", "101", map[string]interface{}{"email": "j@k.l"})
}

func TestCodeWriterLine(t *testing.T) {
	w := &testWriter{records: []logger.Record{}}
	cw := logger.NewCodeWriter(w, logger.LevelDebug)
	cw.SetFileLineLevel("github.com/go-msvc/logger_test/code-writer_test.go", 17, logger.LevelInfo) //will not log debug
	l := logger.New().WithLevel(logger.LevelDebug)
	l.SetWriter(cw)
	test(t, w, l)

	funcName := "test" //where it was writter
	w.assert(t, 0, funcName, l.Name(), "DEBUG", "123", map[string]interface{}{})
	w.assert(t, 1, funcName, l.Name(), "INFO", "123", map[string]interface{}{})
	w.assert(t, 2, funcName, l.Name(), "ERROR", "123", map[string]interface{}{})
	w.assert(t, 3, funcName, l.Name(), "DEBUG", "456", map[string]interface{}{"email": "a@b.c"})
	//w.assert(t, 4, funcName, l.Name(), "DEBUG", "789", map[string]interface{}{"email": "j@k.l"})
	w.assert(t, 4, funcName, l.Name(), "DEBUG", "101", map[string]interface{}{"email": "j@k.l"})
}
