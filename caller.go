package logger

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
)

type ICaller interface {
	fmt.Stringer
	Package() string
	Function() string
	File() string
	Line() int
}

//e.g.:
type Caller struct {
	file       string
	line       int
	pkgDotFunc string
}

func GetCaller(skip int) Caller {
	c := Caller{
		file:       "",
		line:       -1,
		pkgDotFunc: "",
	}

	var pc uintptr
	var ok bool
	if pc, c.file, c.line, ok = runtime.Caller(skip); !ok {
		return c
	}

	if fn := runtime.FuncForPC(pc); fn != nil {
		c.pkgDotFunc = fn.Name()
	}
	return c
} //GetCaller()

func (c Caller) String() string {
	return fmt.Sprintf("%s(%d)", path.Base(c.file), c.line)
}

//with Function: "github.com/go-msvc/ms_test.TestCaller"
//return "github.com/go-msvc/ms_test"
func (c Caller) Package() string {
	if i := strings.LastIndex(c.pkgDotFunc, "."); i >= 0 {
		return c.pkgDotFunc[:i]
	}
	return ""
}

//with Function: "github.com/go-msvc/ms_test.TestCaller"
//return "github.com/go-msvc/ms_test"
func (c Caller) Function() string {
	if i := strings.LastIndex(c.pkgDotFunc, "."); i >= 0 {
		return c.pkgDotFunc[i+1:]
	}
	return ""
}

func (c Caller) File() string {
	return c.file
}

func (c Caller) Line() int {
	return c.line
}

func (caller Caller) Format(f fmt.State, c rune) {
	var s string
	switch c {
	case 's', 'v':
		if p, ok := f.Precision(); ok {
			s = fmt.Sprintf("%s(%*d)", path.Base(caller.file), p, caller.line)
		} else {
			s = fmt.Sprintf("%s(%d)", path.Base(caller.file), caller.line)
		}
	case 'S', 'V':
		if p, ok := f.Precision(); ok {
			s = fmt.Sprintf("%s/%s(%*d)", caller.Package(), path.Base(caller.file), p, caller.line)
		} else {
			s = fmt.Sprintf("%s/%s(%d)", caller.Package(), path.Base(caller.file), caller.line)
		}
	case 'f':
		if p, ok := f.Precision(); ok {
			s = fmt.Sprintf("%s(%*d)", caller.Function(), p, caller.line)
		} else {
			s = fmt.Sprintf("%s(%d)", caller.Function(), caller.line)
		}
	case 'F':
		if p, ok := f.Precision(); ok {
			s = fmt.Sprintf("%s.%s(%*d)", caller.Package(), caller.Function(), p, caller.line)
		} else {
			s = fmt.Sprintf("%s.%s(%d)", caller.Package(), caller.Function(), caller.line)
		}
	} //switch(c)

	if w, ok := f.Width(); ok {
		if w < 0 {
			s = ""
		} else if w <= len(s) {
			//right/left truncate
			if f.Flag('-') {
				s = s[:w]
			} else {
				s = s[len(s)-w:]
			}
		} else if w > len(s) {
			//pad to fill the space
			s = fmt.Sprintf("%*s", w, s)
		}
	}
	io.WriteString(f, s)
} // caller.Format()

//funcName removes the path prefix component of a function's name reported by func.Name().
func funcName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
} // funcName()
