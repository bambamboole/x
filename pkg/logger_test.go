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
			debug: DebugOff,
			out:   &b,
		}

		l.Log("test")

		assert.Equal(t, "test\n", b.String())
	})
	t.Run("It can log properties of structs", func(t *testing.T) {
		var b bytes.Buffer
		l := &IOLogger{
			debug: DebugOff,
			out:   &b,
		}
		testStruct := struct {
			Name string
		}{Name: "test"}

		l.Log(testStruct)

		assert.Equal(t, "name: test\n\n", b.String())
	})
}
