package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunPython(t *testing.T) {
	name := "test"
	data := ""
	want := ""
	var ret string
	err := runPython(name, data, &ret)
	assert.Nil(t, err)
	assert.Equal(t, ret, want, "returned value should be equal")
}
