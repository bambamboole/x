package pkg

import (
	"os"
	"path"
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

func DetectProjectPath(startPath string) string {
	for {
		if isGitRepo(startPath) {
			return startPath
		}
		if isRoot(startPath) {
			return ""
		}
		startPath = path.Dir(startPath)
	}
}
