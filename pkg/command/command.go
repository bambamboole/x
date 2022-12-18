package command

import (
	"os/exec"
	"strings"
	"x/pkg/args"
	"x/pkg/config"
	"x/pkg/utils"
)

type executorInterface interface {
	execute(command string, args ...string) error
}

type executor struct {
	logger utils.Logger
}

func (e *executor) execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	e.logger.Debug("Executing command: " + cmd.String())
	output, err := cmd.CombinedOutput()
	e.logger.Log(string(output))
	return err
}

type Command struct {
	args     args.Arguments
	config   config.Config
	executor executorInterface
	logger   utils.Logger
}

func (c *Command) Execute() error {
	firstArg := c.args.Command[0]
	leftoverArgs := c.args.Command[1:]
	if executable, found := c.config.Executables[firstArg]; found {
		return c.executor.execute(executable.Path, leftoverArgs...)
	}
	bash, _ := exec.LookPath("bash")
	leftoverArgs = append([]string{"source Taskfile", "&&", "task:" + firstArg}, leftoverArgs...)
	return c.executor.execute(bash, "-c", strings.Join(leftoverArgs, " "))
}

func New(args args.Arguments, cfg config.Config, logger utils.Logger) (*Command, error) {
	cmd := &Command{
		args:     args,
		config:   cfg,
		executor: &executor{logger: logger},
		logger:   logger,
	}

	return cmd, nil
}
