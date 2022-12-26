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
