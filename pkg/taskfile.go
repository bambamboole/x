package pkg

import (
	"errors"
	"os"
	"path"
	"strings"
)

type Taskfile struct {
	path   string
	script string
}

func findTaskFilePath(searchPath string) (string, error) {
	run := true

	taskfilePath := ""

	for run {
		if fileExists(searchPath + "/Taskfile") {
			taskfilePath = strings.Clone(searchPath + "/Taskfile")
			run = false
		}

		if folderExists(searchPath + "/.git") {
			run = false
		}
		if searchPath == "/" {
			run = false
		}
		searchPath = path.Dir(searchPath)
	}

	if taskfilePath == "" {
		return taskfilePath, errors.New("Taskfile not found.")
	}

	return taskfilePath, nil
}

func NewTaskfile(searchPath string) (Taskfile, error) {
	taskFilePath, err := findTaskFilePath(searchPath)
	if err != nil {
		return Taskfile{}, err
	}
	script, err := os.ReadFile(taskFilePath)
	if err != nil {
		return Taskfile{}, err
	}
	tf := Taskfile{path: taskFilePath, script: string(script)}

	return tf, nil
}
