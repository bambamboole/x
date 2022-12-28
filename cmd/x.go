package cmd

import (
	x "github.com/bambamboole/x/pkg"
	"github.com/jessevdk/go-flags"
	"os"
)

func Execute() {
	stdout := os.Stdout
	arguments, err := x.ParseArgs(os.Args[1:], stdout)
	if err != nil {
		if flags.WroteHelp(err) {
			return
		}
		panic(err)
	}
	logger := x.NewLogger(len(arguments.Verbose), stdout)
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
	tf, err := x.NewTaskfile(logger, projectPath, arguments.Taskfiles, os.ReadFile)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	cfg, err := x.NewConfig(logger, projectPath, arguments.ConfigFiles)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	commandManager := &x.CommandManager{
		Stdin:  os.Stdin,
		Stdout: stdout,
		Stderr: os.Stderr,
	}

	mr := x.NewModuleRegistry()
	runtime, err := x.NewRuntime(mr, commandManager, projectPath, arguments, cfg, tf, logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	err = runtime.Run()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
