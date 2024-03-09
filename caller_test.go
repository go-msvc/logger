package logger_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/go-msvc/logger"
)

func TestCaller(t *testing.T) {
	lineNr := 13
	c := logger.GetCaller(1)
	t.Logf("Pkg=%s, File=%s, Line=%d, Func=%s", c.Package(), c.File(), c.Line(), c.Function())
	if c.Package() != "github.com/go-msvc/logger_test" {
		t.Fatalf("Package=%s != github.com/go-msvc/logger_test", c.Package())
	}
	if path.Base(c.File()) != "caller_test.go" {
		t.Fatalf("Package=%s != caller_test.go", c.File())
	}
	if c.Line() != lineNr {
		t.Fatalf("Line=%d != %d", c.Line(), lineNr)
	}
	if c.Function() != "TestCaller" {
		t.Fatalf("Function=%s != TestCaller", c.Function())
	}

	//test formatting:
	tests := []struct {
		format        string
		expectedValue string
	}{
		{"%s", fmt.Sprintf("caller_test.go(%d)", lineNr)},
		{"%.5s", fmt.Sprintf("caller_test.go(%5d)", lineNr)},
		{"%S", fmt.Sprintf("github.com/go-msvc/logger_test/caller_test.go(%d)", lineNr)},
		{"%.5S", fmt.Sprintf("github.com/go-msvc/logger_test/caller_test.go(%5d)", lineNr)},
		{"%f", fmt.Sprintf("TestCaller(%d)", lineNr)},
		{"%.5f", fmt.Sprintf("TestCaller(%5d)", lineNr)},
		{"%F", fmt.Sprintf("github.com/go-msvc/logger_test.TestCaller(%d)", lineNr)},
		{"%.5F", fmt.Sprintf("github.com/go-msvc/logger_test.TestCaller(%5d)", lineNr)},
		//precision truncates from left/right
		{"%-10.5F", "github.com"},
		{"%10.5F", fmt.Sprintf("ler(%5d)", lineNr)},
	}
	for index, test := range tests {
		s := fmt.Sprintf(test.format, c)
		if s != test.expectedValue {
			t.Fatalf("test[%d] fmt.Sprintf(\"%s\", caller) -> \"%s\" != \"%s\"", index, test.format, s, test.expectedValue)
		}
		t.Logf("test[%d] OK: fmt.Sprintf(\"%s\", caller) -> \"%s\"", index, test.format, s)
	}
}
