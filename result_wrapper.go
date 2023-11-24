package gocommandinvoker

import (
	"os"
	"os/exec"
)

// newErrResult 创建错误结果
func newErrResult(err error) *Result {
	return &Result{err: err}
}

// newResult 创建结果
func newResult(cmd *exec.Cmd) *Result {
	r := &Result{}

	go r.run()

	return r
}

// Result 命令调用结果
type Result struct {
	cmd     *exec.Cmd
	err     error
	tempDir string
	stdIn   *os.File
	stdOut  *os.File
	stdErr  *os.File
}

func (r *Result) IsError() bool {
	return r.err != nil
}

func (r *Result) Error() error {
	return r.err
}

func (r *Result) run() {
	tempDir, err := os.MkdirTemp("", "gocommandinvoker*")
	if err != nil {
		r.err = ErrCreateTempDirFailed
		return
	}

}
