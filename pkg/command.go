package pkg

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type executorInterface interface {
	execute(command string, args ...string) error
}

type executor struct {
	logger Logger
}

func (e *executor) captureOutput(stdout io.ReadCloser, stderr io.ReadCloser) {
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			e.logger.Log(in.Text())
		}
	}()
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			e.logger.Log("Error: " + in.Text())
		}
	}()
}

func (e *executor) execute(command string, args ...string) error {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGXFSZ)
	cmd := exec.Command(command, args...)

	e.logger.Debug("Executing command: " + cmd.String())
	go func() {
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		_ = cmd.Start()
		e.captureOutput(stdout, stderr)
		_ = cmd.Wait()
		e.logger.Log("Finished")
		cancelChan <- syscall.SIGXFSZ
	}()
	e.logger.Log("Waiting...")

	sig := <-cancelChan
	if sig == syscall.SIGXFSZ {
		return nil
	}
	e.logger.Debug("Got signal: " + sig.String())
	e.logger.Debug("Forwarding cancellation to process...")
	return cmd.Process.Signal(sig)
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
	if executable, found := c.config.Executables[firstArg]; found {
		return c.executor.execute(executable.Path, c.args.Command[1:]...)
	}
	bash, _ := exec.LookPath("bash")
	return c.executor.execute(bash, "-c", c.Taskfile.script+"\ntask:"+strings.Join(c.args.Command, " "))
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
