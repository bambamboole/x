package main

import (
	"os"
	"x/pkg/args"
	"x/pkg/command"
	"x/pkg/config"
	"x/pkg/utils"
)

func main() {
	arguments, err := args.ParseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg, err := config.New(cwd, arguments.ConfigFiles)
	if err != nil {
		panic(err)
	}
	logger := utils.NewLogger(arguments.Verbose)

	cmd, err := command.New(arguments, cfg, logger)
	if err != nil {
		panic(err)
	}
	err = cmd.Execute()
	if err != nil {
		panic(err)
	}
}
