package gocommandinvoker

import (
	"fmt"
	"testing"
)

func TestProcessor_Run(t *testing.T) {

	runner := New()
	r := runner.ExecWithOptions("ls", &RunnerOptions{
		Dir: "~/Desktop",
	}).Run()
	fmt.Println(r)

}
