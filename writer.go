package logger

import (
	"fmt"
	"os"
)

type IWriter interface {
	Write(Record)
}

type defaultWriter struct{}

func (w defaultWriter) Write(r Record) {
	line := fmt.Sprintf("%s %5.5s %25.5s: %s",
		r.Timestamp.Format("2006-01-02 15:04:05.000"),
		r.Level.String(),
		r.Caller,
		r.Message)
	if len(r.Data) > 0 {
		line += fmt.Sprintf(" %v", r.Data)
	}
	fmt.Fprintf(os.Stderr, line+"\n")
}
