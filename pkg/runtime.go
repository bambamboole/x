package pkg

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Runtime struct {
	registry       *ModuleRegistry
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
	script := `
# Initialize error handling
set -o errexit
set -o errtrace
set -o pipefail

# Font colors
X_FONT_GREEN='\033[00;32m'
X_FONT_YELLOW='\033[00;33m'
X_FONT_RED='\033[00;31m'
X_FONT_RESTORE='\033[0m'
X_BASE_PATH='` + r.projectPath + `'

x:warn() {
    echo -e "${X_FONT_YELLOW}${1:-}${X_FONT_RESTORE}" >&2
}

x:error() {
    echo -e "${X_FONT_RED}${1:-}${X_FONT_RESTORE}" >&2
}

`
	for moduleName, module := range r.registry.GetModules() {
		err := r.config.PopulateModuleConfig(moduleName, module.GetConfig())
		if err != nil {
			r.logger.Error("Error while unmarshalling 'modules." + moduleName + "' from config")
			return err
		}
		r.logger.Log(moduleName+" module config:", DebugVeryVerbose)
		r.logger.Log(module.GetConfig(), DebugVeryVerbose)
		moduleScript, err := module.GetScript()
		if err != nil {
			r.logger.Error("Error while building script for module '" + moduleName + "'")
			return err
		}
		script = script + "\n" + moduleScript
	}
	script = strings.TrimPrefix(script, "\n")
	r.logger.Log("Generated module scripts: \n"+script, DebugVeryVerbose)
	r.logger.Log("Generated Taskfile scripts: \n"+r.Taskfile.script, DebugVerbose)
	script = script + "\n" + r.Taskfile.script
	task := "task:" + strings.Join(r.args.Command, " ")
	r.logger.Log("Executing command: " + task)
	script = script + "\n" + task
	return r.execute(shell, "-c", script)
}

func NewRuntime(
	registry *ModuleRegistry,
	commandManager CommandManagerInterface,
	projectPath string,
	args Arguments,
	cfg Config,
	tf Taskfile,
	logger IOLoggerInterface,
) (*Runtime, error) {
	cmd := &Runtime{
		registry:       registry,
		commandManager: commandManager,
		projectPath:    projectPath,
		args:           args,
		config:         cfg,
		Taskfile:       tf,
		logger:         logger,
	}

	return cmd, nil
}
