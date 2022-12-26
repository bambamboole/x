package pkg

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type executorInterface interface {
	execute(workingDir string, command string, args ...string) error
}

type executor struct {
	logger IOLoggerInterface
}

func (e *executor) execute(workingDir string, command string, args ...string) error {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGXFSZ)
	cmd := exec.Command(command, args...)
	e.logger.Log(workingDir)
	cmd.Dir = workingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		_ = cmd.Run()
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
	projectPath string
	cwd         string
	args        Arguments
	config      Config
	Taskfile    Taskfile
	executor    executorInterface
	logger      IOLoggerInterface
}

func (r *Runtime) Execute() error {
	firstArg := r.args.Command[0]
	if executable, found := r.config.Executables[firstArg]; found {
		return r.executor.execute(r.projectPath, executable.Path, r.args.Command[1:]...)
	}
	bash, _ := exec.LookPath(r.args.Shell)
	task := "task:" + strings.Join(r.args.Command, " ")
	r.logger.Log("Using Taskfile content: \n"+r.Taskfile.script, DebugVerbose)
	r.logger.Log("Executing command: "+task, DebugOn)
	return r.executor.execute(r.projectPath, bash, "-c", r.Taskfile.script+"\n"+task)
}

func NewRuntime(projectPath string, cwd string, args Arguments, cfg Config, tf Taskfile, logger IOLoggerInterface) (*Runtime, error) {
	cmd := &Runtime{
		projectPath: projectPath,
		cwd:         cwd,
		args:        args,
		config:      cfg,
		Taskfile:    tf,
		executor:    &executor{logger: logger},
		logger:      logger,
	}

	return cmd, nil
}
