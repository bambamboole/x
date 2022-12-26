package pkg

import (
	"github.com/jessevdk/go-flags"
	"os"
)

type Arguments struct {
	Verbose     []bool   `short:"v" long:"verbose" description:"Enable debug output"`
	ConfigFiles []string `short:"c" long:"config-files" description:"Additional config files to load" default:"~/.x/config.yml"`
	Taskfiles   []string `short:"t" long:"taskfiles" description:"Additional task files to load" default:"~/.x/Taskfile"`
	Shell       string   `short:"s" long:"shell" description:"The shell to execute the commands" default:"bash"`
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
