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
	fmt.Fprintf(os.Stderr, "%s %5.5s %25.5s: %s %+v\n",
		r.Timestamp.Format("2006-01-02 15:04:05.000"),
		r.Level.String(),
		r.Caller,
		r.Message,
		r.Data,
	)
}
