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
	logger IOLoggerInterface
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

	go func() {
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		_ = cmd.Start()
		e.captureOutput(stdout, stderr)
		_ = cmd.Wait()
		cancelChan <- syscall.SIGXFSZ
	}()

	sig := <-cancelChan
	if sig == syscall.SIGXFSZ {
		return nil
	}
	e.logger.Log("Got signal: "+sig.String(), DebugOn)
	e.logger.Log("Forwarding cancellation to process...", DebugOn)
	return cmd.Process.Signal(sig)
}

type Runtime struct {
	args     Arguments
	config   Config
	Taskfile Taskfile
	executor executorInterface
	logger   IOLoggerInterface
}

func (r *Runtime) Execute() error {
	firstArg := r.args.Command[0]
	if executable, found := r.config.Executables[firstArg]; found {
		return r.executor.execute(executable.Path, r.args.Command[1:]...)
	}
	bash, _ := exec.LookPath(r.args.Shell)
	task := "task:" + strings.Join(r.args.Command, " ")
	r.logger.Log("Using Taskfile content: \n"+r.Taskfile.script, DebugVerbose)
	r.logger.Log("Executing command: "+task, DebugOn)
	return r.executor.execute(bash, "-c", r.Taskfile.script+"\n"+task)
}

func NewRuntime(args Arguments, cfg Config, tf Taskfile, logger IOLoggerInterface) (*Runtime, error) {
	cmd := &Runtime{
		args:     args,
		config:   cfg,
		Taskfile: tf,
		executor: &executor{logger: logger},
		logger:   logger,
	}

	return cmd, nil
}
