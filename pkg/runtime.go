package pkg

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type Runtime struct {
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	projectPath string
	args        Arguments
	config      Config
	Taskfile    Taskfile
	logger      IOLoggerInterface
}

func (r *Runtime) execute(command string, args ...string) error {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGXFSZ)
	cmd := exec.Command(command, args...)
	cmd.Dir = r.projectPath
	cmd.Stdin = r.stdin
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr

	go func() {
		_ = cmd.Run()
		cancelChan <- syscall.SIGXFSZ
	}()

	sig := <-cancelChan
	if sig == syscall.SIGXFSZ {
		return nil
	}
	r.logger.Log("Got signal: "+sig.String(), DebugOn)
	r.logger.Log("Forwarding cancellation to process...", DebugOn)
	return cmd.Process.Signal(sig)
}

func (r *Runtime) Run() error {
	shell, _ := exec.LookPath(r.args.Shell)
	task := "task:" + strings.Join(r.args.Command, " ")
	r.logger.Log("Using Taskfile content: \n"+r.Taskfile.script, DebugVerbose)
	r.logger.Log("Executing command: " + task)
	return r.execute(shell, "-c", r.Taskfile.script+"\n"+task)
}

func NewRuntime(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	projectPath string,
	args Arguments,
	cfg Config,
	tf Taskfile,
	logger IOLoggerInterface,
) (*Runtime, error) {
	cmd := &Runtime{
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		projectPath: projectPath,
		args:        args,
		config:      cfg,
		Taskfile:    tf,
		logger:      logger,
	}

	return cmd, nil
}
