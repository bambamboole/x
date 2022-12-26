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
		logger.Error(err)
		return
	}
	projectPath := x.DetectProjectPath(cwd)
	logger.Log("Use project path: "+projectPath, x.DebugOn)
	tf, err := x.NewTaskfile(logger, projectPath, arguments.Taskfiles)
	if err != nil {
		logger.Error(err)
		return
	}
	cfg, err := x.NewConfig(logger, projectPath, arguments.ConfigFiles)
	if err != nil {
		logger.Error(err)
		return
	}

	runtime, err := x.NewRuntime(projectPath, cwd, arguments, cfg, tf, logger)
	if err != nil {
		logger.Error(err)
		return
	}
	err = runtime.Execute()
	if err != nil {
		logger.Error(err)
		return
	}
}
