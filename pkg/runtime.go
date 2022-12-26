package pkg

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Runtime struct {
	commandManager CommandManagerInterface
	projectPath    string
	args           Arguments
	config         Config
	Taskfile       Taskfile
	logger         IOLoggerInterface
}

func (r *Runtime) execute(command string, args ...string) error {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGXFSZ)
	cmd := r.commandManager.Create(r.projectPath, command, args)

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
	return r.commandManager.Stop(cmd, sig)
}

func (r *Runtime) Run() error {
	shell, _ := r.commandManager.FindExecutable(r.args.Shell)
	task := "task:" + strings.Join(r.args.Command, " ")
	r.logger.Log("Using Taskfile content: \n"+r.Taskfile.script, DebugVerbose)
	r.logger.Log("Executing command: " + task)
	return r.execute(shell, "-c", r.Taskfile.script+"\n"+task)
}

func NewRuntime(
	commandManager CommandManagerInterface,
	projectPath string,
	args Arguments,
	cfg Config,
	tf Taskfile,
	logger IOLoggerInterface,
) (*Runtime, error) {
	cmd := &Runtime{
		commandManager: commandManager,
		projectPath:    projectPath,
		args:           args,
		config:         cfg,
		Taskfile:       tf,
		logger:         logger,
	}

	return cmd, nil
}
