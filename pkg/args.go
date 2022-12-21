package pkg

import (
	"fmt"
	"github.com/jessevdk/go-flags"
)

type Arguments struct {
	Verbose     bool     `short:"v" long:"verbose" description:"Enable debug output"`
	ConfigFiles []string `short:"c" long:"config-files" description:"Additional config files to be loaded" default:"~/.x/config.yml"`
	Command     []string
}

func ParseArgs(args []string) (Arguments, error) {
	a := Arguments{}
	p := flags.NewParser(&a, flags.IgnoreUnknown)
	args, err := p.Parse()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return a, err
	}
	a.Command = args

	return a, nil
}
