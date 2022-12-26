package pkg

import (
	"github.com/jessevdk/go-flags"
	"io"
	"strings"
)

type Arguments struct {
	Verbose     []bool   `short:"v" long:"verbose" description:"Enable debug output"`
	ConfigFiles []string `short:"c" long:"config-files" description:"Additional config files to load" default:"~/.x/config.yml"`
	Taskfiles   []string `short:"t" long:"taskfiles" description:"Additional task files to load" default:"~/.x/Taskfile"`
	Shell       string   `short:"s" long:"shell" description:"The shell to execute the commands" default:"bash"`
	Command     []string
}

func ParseArgs(args []string, stdout io.Writer) (Arguments, error) {
	a := Arguments{}
	p := flags.NewParser(&a, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
	args, err := p.ParseArgs(args)
	if len(args) == 0 {
		err = &flags.Error{Type: flags.ErrHelp}
	}
	if err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type == flags.ErrHelp {
			p.WriteHelp(stdout)

		}
		return a, err
	}
	// append verbose flag if present, since it gets removed by the flags parser
	if len(a.Verbose) > 0 {
		args = append(args, "-"+strings.Repeat("v", len(a.Verbose)))
	}
	a.Command = args

	return a, nil
}
