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

func TestInvalidTimeout(t *testing.T) {
	args := []string{
		"test/echo_time.sh"}

	err := scripts(args, -1)
	if err == nil {
		t.Errorf("not rejected invalid timeout :%s\n", err)
	}
}

func Test1ScriptTimeout(t *testing.T) {
	args := []string{
		"test/echo_time.sh"}

	err := scripts(args, 1)
	if err != nil {
		t.Errorf("Not succeeded to exec2scripts:%s\n", err)
	}
}

func Test1ScriptFinish(t *testing.T) {
	args := []string{
		"test/echo_time_1sec.sh"}

	err := scripts(args, 10)
	if err != nil {
		t.Errorf("Not succeeded to exec2scripts:%s\n", err)
	}
}

func TestScriptsTimeout(t *testing.T) {
	args := []string{
		"test/echo_time.sh",
		"test/echo_cnt.sh"}

	err := scripts(args, 3)
	if err != nil {
		t.Errorf("Not succeeded to exec2scripts:%s\n", err)
	}
}

func TestScriptsFinish(t *testing.T) {
	args := []string{
		"test/echo_time_1sec.sh",
		"test/echo_cnt_2sec.sh"}

	err := scripts(args, 3)
	if err != nil {
		t.Errorf("Not succeeded to exec2scripts:%s\n", err)
	}
}

func TestLeak(t *testing.T) {
	defer goleak.VerifyNone(t)

	// test logic here.
}
