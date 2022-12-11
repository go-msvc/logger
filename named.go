package logger

import (
	"sync"
)

type named struct {
	name  string
	names []string

	level Level

	sync.Mutex
	parent *named
	subs   map[string]*named

	writer IWriter
}

func (l *named) New(name string) Logger {
	if l == nil {
		return top.New(name)
	}

	l.Lock()
	defer l.Unlock()
	nl, found := l.subs[name]
	if !found {
		nl = &named{
			name:   name,
			names:  append(l.names, name),
			parent: l,
			subs:   map[string]*named{},
			level:  l.level,
			writer: l.writer,
		}
		l.subs[name] = nl
	}

	return logger{
		named: nl,
		data:  map[string]interface{}{},
		level: LevelDefault,
	}
} //named.New()

func (l *named) setWriter(newWriter IWriter) {
	if newWriter == nil {
		newWriter = defaultWriter{}
	}

	l.Lock()
	defer l.Unlock()
	for _, sub := range l.subs {
		sub.setWriter(newWriter)
	}
	l.writer = newWriter
}

func (l *named) setLevel(newLevel Level) {
	l.Lock()
	defer l.Unlock()
	for _, sub := range l.subs {
		sub.setLevel(newLevel)
	}
	l.level = newLevel
}

func (l *named) Name() string    { return l.name }
func (l *named) Names() []string { return l.names }

func (l *named) Level() Level { return l.level }

func (l *named) WithLevel(newLevel Level) Logger {
	return logger{
		named: l,
		level: LevelDefault,
		data:  map[string]interface{}{},
	}
}

func (l *named) With(name string, value interface{}) Logger {
	return logger{
		named: l,
		level: LevelDefault,
		data:  map[string]interface{}{name: value},
	}
}

type INamed interface {
	Name() string
	Names() []string
	Level() Level
	All() map[string]INamed
}

func (l *named) All() map[string]INamed {
	all := map[string]INamed{}
	for _, s := range l.subs {
		all[s.name] = s
	}
	return all
}
