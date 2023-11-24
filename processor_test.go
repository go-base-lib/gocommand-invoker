package gocommandinvoker

import (
	"fmt"
	"testing"
)

func TestProcessor_Run(t *testing.T) {

	runner := New()
	r := runner.ExecWithOptions("ls", &RunnerOptions{
		Dir: "~",
	}).Run()
	fmt.Println(r)

}
