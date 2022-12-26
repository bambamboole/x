package main

import (
	x "github.com/bambamboole/x/pkg"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	arguments, err := x.ParseArgs(os.Args[1:])
	if err != nil {
		if flags.WroteHelp(err) {
			return
		}
		panic(err)
	}
	logger := x.NewLogger(len(arguments.Verbose))
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg, err := x.NewConfig(cwd, arguments.ConfigFiles)
	if err != nil {
		panic(err)
	}
	tf, err := x.NewTaskfile(cwd)
	if err != nil {
		panic(err)
	}
	runtime, err := x.NewRuntime(arguments, cfg, tf, logger)
	if err != nil {
		panic(err)
	}
	err = runtime.Execute()
	if err != nil {
		panic(err)
	}
}
