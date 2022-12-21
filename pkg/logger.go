package pkg

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Log(v any)
	Debug(v any)
}

type IOLogger struct {
	debug bool
	out   io.Writer
}

func (l *IOLogger) Log(v any) {
	switch v.(type) {
	default:
		_, _ = fmt.Fprintf(l.out, "%#v\n", v)
	case string:
		_, _ = fmt.Fprintln(l.out, v)
	}
}

func (l *IOLogger) Debug(v any) {
	if l.debug == false {
		return
	}
	l.Log(v)
}

func NewLogger(debug bool) Logger {
	logger := &IOLogger{
		debug: debug,
		out:   os.Stdout,
	}
	logger.Debug("Debug output enabled")

	return logger
}
