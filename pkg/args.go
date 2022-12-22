package pkg

import (
	"github.com/jessevdk/go-flags"
	"os"
)

type Arguments struct {
	Verbose     bool     `short:"v" long:"verbose" description:"Enable debug output"`
	ConfigFiles []string `short:"c" long:"config-files" description:"Additional config files to be loaded" default:"~/.x/config.yml"`
	Command     []string
}

func ParseArgs(args []string) (Arguments, error) {
	a := Arguments{}
	p := flags.NewParser(&a, flags.IgnoreUnknown|flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	args, err := p.Parse()
	if err != nil {
		return a, err
	}
	if len(args) == 0 {
		p.WriteHelp(os.Stdout)
		return a, &flags.Error{Type: flags.ErrHelp}
	}
	a.Command = args

	return a, nil
}
