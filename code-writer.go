package logger

type ICodeWriter interface {
	IWriter
	SetFileLevel(name string, level Level)               //level=Default to delete all file settings
	SetFileLineLevel(name string, line int, level Level) //level=Default to delete line setting
}

func NewCodeWriter(w IWriter, l Level) ICodeWriter {
	return &codeWriter{
		writer:    w,
		level:     l,
		fileLevel: map[string]fileLevel{},
	}
}

type codeWriter struct {
	writer    IWriter
	level     Level
	fileLevel map[string]fileLevel
}

type fileLevel struct {
	level     Level
	lineLevel map[int]Level
}

func (cw codeWriter) Write(record Record) {
	//see if should write
	fileName := record.Caller.PackageFile()
	fileLevel, ok := cw.fileLevel[fileName]
	if !ok {
		//no file entry - use global level
		if record.Level > cw.level {
			return //do not write global
		}
	} else {
		//has file entry
		lineLevel, ok := fileLevel.lineLevel[record.Caller.Line()]
		if !ok {
			//no line entry
			if record.Level > fileLevel.level {
				return //do not write this file
			}
		} else {
			if record.Level > lineLevel {
				//file.line has an entry
				return //do not write this line
			}
		}
	}
	cw.writer.Write(record)
}

func (cw *codeWriter) SetFileLevel(name string, level Level) {
	if name != "" {
		if level == LevelDefault {
			delete(cw.fileLevel, name)
		} else {
			fl, ok := cw.fileLevel[name]
			if !ok {
				fl = fileLevel{level: level, lineLevel: map[int]Level{}}
			}
			cw.fileLevel[name] = fl
		}
	}
}

func (cw *codeWriter) SetFileLineLevel(name string, line int, level Level) {
	if name != "" && line > 0 {
		fl, ok := cw.fileLevel[name]
		if !ok {
			fl = fileLevel{level: LevelDefault, lineLevel: map[int]Level{}}
		}
		if level == LevelDefault {
			delete(fl.lineLevel, line)
		} else {
			fl.lineLevel[line] = level
		}
		if len(fl.lineLevel) == 0 && fl.level == LevelDefault {
			delete(cw.fileLevel, name) //nothing remain for this file, get rid of it
		} else {
			cw.fileLevel[name] = fl
		}
	}
}
