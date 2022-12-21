package pkg

import (
	"os/exec"
	"strings"
)

type executorInterface interface {
	execute(command string, args ...string) error
}

type executor struct {
	logger Logger
}

func (e *executor) execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	e.logger.Debug("Executing command: " + cmd.String())
	output, err := cmd.CombinedOutput()
	e.logger.Log(string(output))
	return err
}

type Command struct {
	args     Arguments
	config   Config
	Taskfile Taskfile
	executor executorInterface
	logger   Logger
}

func (c *Command) Execute() error {
	firstArg := c.args.Command[0]
	leftoverArgs := c.args.Command[1:]
	if executable, found := c.config.Executables[firstArg]; found {
		return c.executor.execute(executable.Path, leftoverArgs...)
	}
	bash, _ := exec.LookPath("bash")
	leftoverArgs = append([]string{"source", c.Taskfile.path, "&&", "task:" + firstArg}, leftoverArgs...)
	return c.executor.execute(bash, "-c", strings.Join(leftoverArgs, " "))
}

func NewCommand(args Arguments, cfg Config, tf Taskfile, logger Logger) (*Command, error) {
	cmd := &Command{
		args:     args,
		config:   cfg,
		Taskfile: tf,
		executor: &executor{logger: logger},
		logger:   logger,
	}

	return cmd, nil
}
