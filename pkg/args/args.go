package args

import "github.com/jessevdk/go-flags"

type Arguments struct {
	ConfigFiles []string `short:"c" long:"config-files" description:"Additional config files to be loaded" default:"~/.x/config.yml"`
	Command     []string
}

func ParseArgs(args []string) (*Arguments, error) {
	a := &Arguments{}
	args, err := flags.ParseArgs(a, args)
	if err != nil {
		return nil, err
	}
	a.Command = args

	return a, nil
}
