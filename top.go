package logger

func New() Logger {
	caller := GetCaller(2)
	return top.New(caller.Package())
}

func Named(name string) Logger {
	return top.New(name)
}

var (
	top *named
)

func init() {
	top = &named{
		name:   "",
		parent: nil,
		subs:   map[string]*named{},
		level:  LevelError,
		writer: defaultWriter{},
	}
}

// SetGlobalWriter sets the writer on top and all existing and default for all new loggers
func SetGlobalWriter(newWriter IWriter) {
	top.setWriter(newWriter)
}

// careful with this - will also set for all named loggers used in libraries
// rather create your own logger and set its level, or set level on Named(...).SetLevel(...)
// still: setting level like this will only affect loggers that did not set their own levels
// libraries should just create logger that defaults to level error, and then will be changed
// by this or Named(...).SetLevel()
// If a library did New().WithLevel(), it will not be affected by any of this...
// so libraries should not set their levels! Then level is in the named.level and all loggers
// created for it will use LevelDefault to inherit from the named...
func SetGlobalLevel(newLevel Level) {
	top.setLevel(newLevel)
}

func All() map[string]INamed {
	return top.All()
}
