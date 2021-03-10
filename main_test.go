package main

import (
	"testing"

	"go.uber.org/goleak"
)

func TestIsExistsScript(t *testing.T) {
	isExists := isExistsScript("test/echo_time.sh")

	if !isExists {
		t.Error("not found scirpt")
	}
	isExists = isExistsScript("test/no_exist.sh")

	if isExists {
		t.Error("found scirpt wrong")
	}
}

func TestExec2Scripts(t *testing.T) {
	scripts := []string{
		"test/echo_time.sh",
		"test/echo_cnt.sh"}

	err := execScripts(scripts)
	if err != nil {
		t.Errorf("Not succeeded to exec2scripts:%s\n", err)
	}
}

func TestA(t *testing.T) {
	defer goleak.VerifyNone(t)

	// test logic here.
}
