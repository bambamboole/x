package pkg

import (
	"os"
	"path"
	"strings"
)

type Taskfile struct {
	script string
}

func NewTaskfile(logger IOLoggerInterface, projectPath string, additionalTaskFiles []string) (Taskfile, error) {
	taskFiles := append(additionalTaskFiles, path.Join(projectPath, "Taskfile"), path.Join(projectPath, "Taskfile.local"))
	scriptContent := ""
	for _, taskFile := range taskFiles {
		if !fileExists(taskFile) {
			logger.Log("No Taskfile found at: "+taskFile, DebugOn)
			continue
		}
		logger.Log("Found file at: "+taskFile+" - will be attached Taskfile", DebugOn)

		content, err := os.ReadFile(taskFile)
		if err != nil {
			return Taskfile{}, err
		}
		scriptContent = scriptContent + "\n" + string(content)
	}

	return Taskfile{script: strings.TrimPrefix(scriptContent, "\n")}, nil
}
