package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolPointer(t *testing.T) {
	b := Bool("bool", true, "help")
	assert.Equal(t, true, *b)
	*Lookup("bool").Value.(*bool) = false
	assert.Equal(t, false, *b)
}

func TestStringPointer(t *testing.T) {
	b := String("str", "def", "help")
	assert.Equal(t, "def", *b)
	*Lookup("str").Value.(*string) = "new"
	assert.Equal(t, "new", *b)
}
