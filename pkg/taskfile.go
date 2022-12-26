package pkg

import (
	"os"
	"strings"
)

type Taskfile struct {
	script string
}

func NewTaskfile(searchPath string) (Taskfile, error) {
	taskFiles := findFiles(searchPath, []string{"Taskfile", "Taskfile.local"})
	scriptContent := ""
	for _, taskFile := range taskFiles {
		content, err := os.ReadFile(taskFile)
		if err != nil {
			return Taskfile{}, err
		}
		scriptContent = scriptContent + "\n" + string(content)
	}

	tf := Taskfile{script: strings.TrimPrefix(scriptContent, "\n")}

	return tf, nil
}
