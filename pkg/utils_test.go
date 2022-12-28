package pkg

import (
	"github.com/stretchr/testify/assert"
	"path"
	"runtime"
	"testing"
)

func Test_fileExists(t *testing.T) {
	t.Run("it can check if a file exists", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		assert.True(t, fileExists(path.Join(path.Dir(filename), "fixtures", "test_project", "dummy_file.txt")))
	})
	t.Run("it returns false on a folder for fileExists", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		assert.False(t, fileExists(path.Join(path.Dir(filename), "fixtures", "test_project")))
	})
}

func Test_folderExists(t *testing.T) {
	t.Run("it can check if a folder exists", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		assert.True(t, directoryExists(path.Join(path.Dir(filename), "fixtures", "test_project")))
	})
	t.Run("it returns false on a file for directoryExists", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		assert.False(t, directoryExists(path.Join(path.Dir(filename), "fixtures", "test_project", "dummy_file.txt")))
	})
}

func Test_isGitRepo(t *testing.T) {
	t.Run("it can check if a folder is a git repository", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		assert.True(t, isGitRepo(path.Join(path.Dir(filename), "fixtures", "test_project")))
	})
	t.Run("cross check with non git folder", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		assert.False(t, isGitRepo(path.Join(path.Dir(filename), "fixtures")))
	})
}

func Test_isRoot(t *testing.T) {
	t.Run("it checks if the passed string is the root folder", func(t *testing.T) {
		assert.True(t, isRoot("/"))
		assert.False(t, isRoot("/foobar"))
	})
}

func Test_DetectProjectPath(t *testing.T) {
	t.Run("It detects the project path based on git folder", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)

		projectPath := DetectProjectPath(path.Join(path.Dir(filename), "fixtures", "test_project", "subdir", "subdir"))
		_, projectPath = path.Split(projectPath)

		assert.Equal(t, "test_project", projectPath)
	})
}
