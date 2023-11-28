package gocommandinvoker

import (
	"fmt"
	"testing"
	"time"
)

func TestProcessor_Run(t *testing.T) {

	runner := New()
	r := runner.ExecWithOptions("sleep 15s;echo \"hello world\"", &RunnerOptions{
		Dir: "~",
	}).Run()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("执行了")
		_ = r.Kill()
	}()
	fmt.Println(r.String())
	fmt.Println(r.Pid())

}
