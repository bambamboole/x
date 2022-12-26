package pkg

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"sync"
	"syscall"
	"testing"
)

func TestCommandManager_FindExecutable(t *testing.T) {
	t.Run("It can find the absolute path of an executable form the path", func(t *testing.T) {
		cm := &CommandManager{}
		sh, _ := cm.FindExecutable("sh")
		assert.Equal(t, "/bin/sh", sh)
	})
}

func TestCommandManager_Create(t *testing.T) {
	t.Run("It sets the correct properties on the command", func(t *testing.T) {
		var stdin bytes.Reader
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		path := "/test"

		cm := &CommandManager{Stdin: &stdin, Stdout: &stdout, Stderr: &stderr}
		cmd := cm.Create(path, "test", []string{})
		command, _ := cmd.(*exec.Cmd)

		assert.Equal(t, &stdin, command.Stdin)
		assert.Equal(t, &stdout, command.Stdout)
		assert.Equal(t, &stderr, command.Stderr)
		assert.Equal(t, path, command.Dir)
	})
}

func TestCommandManager_Stop(t *testing.T) {
	t.Run("It passes the signel to the command", func(t *testing.T) {
		var stdin bytes.Reader
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		path := "/"

		cm := &CommandManager{Stdin: &stdin, Stdout: &stdout, Stderr: &stderr}
		sh, _ := cm.FindExecutable("sh")
		cmd := cm.Create(path, sh, []string{"-c", "echo foo"})
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			err := cmd.Run()
			assert.Equal(t, nil, err)
			wg.Done()
		}()
		err := cm.Stop(cmd, syscall.SIGKILL)
		assert.Equal(t, nil, err)
		wg.Wait()
		assert.Equal(t, "foo\n", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
