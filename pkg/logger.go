package pkg

import (
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io"
)

const (
	DebugOff         = 0
	DebugOn          = 1
	DebugVerbose     = 2
	DebugVeryVerbose = 3
)

type IOLoggerInterface interface {
	Log(v any, debugLevel ...int)
	Error(v any)
}

type IOLogger struct {
	debug int
	out   io.Writer
}

func (l *IOLogger) Log(v any, debugLevel ...int) {
	debug := DebugOff
	if len(debugLevel) > 0 {
		debug = debugLevel[0]
	}
	if debug > l.debug {
		return
	}
	col := color.New(color.Reset)
	switch debug {
	case DebugOn:
		col = color.New(color.FgCyan)
	case DebugVerbose:
		col = color.New(color.FgBlue)
	case DebugVeryVerbose:
		col = color.New(color.FgGreen)
	}
	l.write(v, col)
}

func (l *IOLogger) Error(v any) {
	err, ok := v.(error)
	if ok {
		v = err.Error()
	}
	l.write(v, color.New(color.FgRed))
}

func (l *IOLogger) write(v any, col *color.Color) {
	switch v.(type) {
	default:
		vJSON, _ := yaml.Marshal(v)
		_, _ = col.Fprintln(l.out, string(vJSON))
	case string:
		_, _ = col.Fprintln(l.out, v)
	}
}

func NewLogger(debug int, stdout io.Writer) IOLoggerInterface {
	logger := &IOLogger{
		debug: debug,
		out:   stdout,
	}
	logger.Log("Debug output enabled", DebugOn)

	return logger
}
