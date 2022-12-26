package cmd

import (
	x "github.com/bambamboole/x/pkg"
	"github.com/jessevdk/go-flags"
	"os"
)

func Execute() {
	arguments, err := x.ParseArgs(os.Args[1:])
	if err != nil {
		if flags.WroteHelp(err) {
			return
		}
		panic(err)
	}
	logger := x.NewLogger(len(arguments.Verbose))
	logger.Log("x " + x.Version + " (Commit: " + x.Commit + ", Build date: " + x.Date + ")")
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	projectPath := x.DetectProjectPath(cwd)
	if projectPath == "" {
		logger.Error("Was not able to detect project path")
		os.Exit(1)
	}
	logger.Log("Use project path: "+projectPath, x.DebugOn)
	tf, err := x.NewTaskfile(logger, projectPath, arguments.Taskfiles)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	cfg, err := x.NewConfig(logger, projectPath, arguments.ConfigFiles)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	runtime, err := x.NewRuntime(projectPath, cwd, arguments, cfg, tf, logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	err = runtime.Execute()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
