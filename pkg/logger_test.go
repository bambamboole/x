package pkg

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIOLogger_Log(t *testing.T) {
	t.Run("It can logs strings and appends new line", func(t *testing.T) {
		var b bytes.Buffer
		l := &IOLogger{
			debug: false,
			out:   &b,
		}

		l.Log("test")

		assert.Equal(t, "test\n", b.String())
	})
	t.Run("It can log structs", func(t *testing.T) {
		var b bytes.Buffer
		l := &IOLogger{
			debug: false,
			out:   &b,
		}
		testStruct := struct {
			Name string
		}{Name: "test"}

		l.Log(testStruct)

		assert.Equal(t, "struct { Name string }{Name:\"test\"}\n", b.String())
	})
}

func TestIOLogger_Debug(t *testing.T) {
	t.Run("It does not log Debug calls if debug is false", func(t *testing.T) {
		var b bytes.Buffer
		l := &IOLogger{
			debug: false,
			out:   &b,
		}

		l.Debug("test")

		assert.Equal(t, "", b.String())
	})
	t.Run("It does log Debug calls if debug is true", func(t *testing.T) {
		var b bytes.Buffer
		l := &IOLogger{
			debug: true,
			out:   &b,
		}

		l.Debug("test")

		assert.Equal(t, "test\n", b.String())
	})
}
