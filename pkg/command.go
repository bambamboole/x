package pkg

import (
	"io"
	"os"
	"os/exec"
)

type CommandInterface interface {
	Run() error
}

type CommandManagerInterface interface {
	FindExecutable(name string) (string, error)
	Create(path string, command string, args []string) CommandInterface
	Stop(cmd CommandInterface, sig os.Signal) error
}

type CommandManager struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (cm *CommandManager) FindExecutable(name string) (string, error) {
	return exec.LookPath(name)
}

func (cm *CommandManager) Create(path string, command string, args []string) CommandInterface {
	cmd := exec.Command(command, args...)
	cmd.Dir = path
	cmd.Stdin = cm.Stdin
	cmd.Stdout = cm.Stdout
	cmd.Stderr = cm.Stderr

	return cmd
}

func (cm *CommandManager) Stop(cmd CommandInterface, sig os.Signal) error {
	command, ok := cmd.(*exec.Cmd)
	if ok && command.Process != nil {
		return command.Process.Signal(sig)
	}
	return nil
}
