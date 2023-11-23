package gocommandinvoker

import (
	"os"
	"syscall"
)

func (r *RunnerOptions) fillDefault(cmd string) {
	if r.SysProcAttr == nil {
		r.SysProcAttr = &syscall.SysProcAttr{}
	}

	r.Env = append(os.Environ(), r.Env...)

	//os.whi
}
