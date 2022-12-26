package pkg

import (
	"os"
	"path"
	"strings"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func directoryExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func isRoot(p string) bool {
	return p == "/"
}

func isGitRepo(p string) bool {
	return directoryExists(path.Join(p, ".git"))
}

func findFiles(searchPath string, fileNames []string) []string {
	stopFunctions := []func(string) bool{isGitRepo, isRoot}
	files := make([]string, 0)

	run := true

	for run {
		for _, fileName := range fileNames {
			configPath := path.Join(searchPath, fileName)
			if fileExists(configPath) {
				files = append(files, strings.Clone(configPath))
				run = false
			}
		}
		for _, stopFunction := range stopFunctions {
			if stopFunction(searchPath) {
				run = false
			}
		}
		searchPath = path.Dir(searchPath)
	}

	return files
}
