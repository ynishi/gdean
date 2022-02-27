package service

import (
	"testing"
)

func TestRunPython(t *testing.T) {
	name := "test"
	data := ""
	want := ""
	var ret string
	err := runPython(name, data, &ret)
	if ret != want || err != nil {
		t.Fatalf(`Ng in %s, err %s, want %s`, ret, err, want)
	}
}
