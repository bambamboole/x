package pkg

import (
	"github.com/fatih/color"
	"io"
	"os"
)

const (
	DebugOff         = 0
	DebugOn          = 1
	DebugVerbose     = 2
	DebugVeryVerbose = 3
)

type IOLoggerInterface interface {
	Log(v any, debugLevel ...int)
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
	switch v.(type) {
	default:
		_, _ = col.Fprintf(l.out, "%#v\n", v)
	case string:
		_, _ = col.Fprintln(l.out, v)
	}
}

func NewLogger(debug int) IOLoggerInterface {
	logger := &IOLogger{
		debug: debug,
		out:   os.Stdout,
	}
	logger.Log("Debug output enabled", DebugOn)

	return logger
}
