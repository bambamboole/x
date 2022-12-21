package main

import (
	"os"
	x "x/pkg"
)

func main() {
	arguments, err := x.ParseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}
	logger := x.NewLogger(arguments.Verbose)
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
	cmd, err := x.NewCommand(arguments, cfg, tf, logger)
	if err != nil {
		panic(err)
	}
	err = cmd.Execute()
	logger.Log(cmd)
	if err != nil {
		panic(err)
	}
}
