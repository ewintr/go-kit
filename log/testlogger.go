package log

type TestLine struct {
	Level   LogLevel
	Message string
	Fields  Fields
}

type TestOut struct {
	Lines []TestLine
}

func NewTestOut() *TestOut {
	return &TestOut{
		Lines: make([]TestLine, 0),
	}
}

func (to *TestOut) Append(tl TestLine) {
	to.Lines = append(to.Lines, tl)
}

func (to *TestOut) Flush() {
	to.Lines = make([]TestLine, 0)
}

type TestLogger struct {
	fields Fields
	level  LogLevel
	out    *TestOut
}

func NewTestLogger(out *TestOut) Logger {
	return &TestLogger{
		fields: make(Fields),
		level:  LevelDebug,
		out:    out,
	}
}

func (tl *TestLogger) SetLogLevel(level LogLevel) {
	tl.level = level
}

func (tl *TestLogger) WithField(key string, value interface{}) Logger {
	return tl.With(Fields{key: value})
}

func (tl *TestLogger) WithErr(err error) Logger {
	return tl.With(Fields{"error": err})
}

func (tl *TestLogger) With(fields Fields) Logger {
	newFields := make(Fields)
	for k, v := range tl.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &TestLogger{
		fields: newFields,
		level:  tl.level,
		out:    tl.out,
	}
}

func (tl *TestLogger) Debug(message string) {
	tl.out.Append(TestLine{
		Level:   LevelDebug,
		Message: message,
		Fields:  tl.fields,
	})

	tl.fields = make(Fields)
}

func (tl *TestLogger) Info(message string) {
	tl.out.Append(TestLine{
		Level:   LevelInfo,
		Message: message,
		Fields:  tl.fields,
	})

	tl.fields = make(Fields)
}

func (tl *TestLogger) Error(message string) {
	tl.out.Append(TestLine{
		Level:   LevelError,
		Message: message,
		Fields:  tl.fields,
	})

	tl.fields = make(Fields)
}
