package gocommandinvoker

import (
	"errors"
	"io"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

// newErrResult 创建错误结果
func newErrResult(err error) *Result {
	return &Result{err: err}
}

// newResult 创建结果
func newResult(cmd *exec.Cmd) *Result {
	r := &Result{
		cmd: cmd,
	}

	return r.run()
}

// ResultStatus 结果状态
type ResultStatus uint8

const (
	// ResultStatusNoRun 未运行
	ResultStatusNoRun ResultStatus = iota
	// ResultStatusRunning 运行中
	ResultStatusRunning
	// ResultStatusRunError 运行错误
	ResultStatusRunError
	// ResultStatusSuccess 成功
	ResultStatusSuccess
	// ResultStatusRunWaitError 等待错误
	ResultStatusRunWaitError
	// ResultStatusOtherError 错误
	ResultStatusOtherError
)

// Result 命令调用结果
type Result struct {
	sync.Mutex
	cmd    *exec.Cmd
	err    error
	stdIn  io.WriteCloser
	stdOut io.ReadCloser
	stdErr io.ReadCloser
	status ResultStatus
	runOk  chan struct{}
}

// IsError 是否错误
func (r *Result) IsError() bool {
	return r.err != nil
}

// Error 获取错误
func (r *Result) Error() error {
	return r.err
}

// run 运行
func (r *Result) run() *Result {
	r.Lock()
	defer r.Unlock()

	if r.runOk != nil {
		return r
	}

	r.runOk = make(chan struct{})

	var (
		err error
	)

	defer func() {
		if err != nil {
			close(r.runOk)
		}
	}()

	if r.status != ResultStatusNoRun {
		return r
	}

	if r.stdIn, err = r.cmd.StdinPipe(); err != nil {
		r.restOtherErrorStatus(err)
		return r
	}

	if r.stdOut, err = r.cmd.StdoutPipe(); err != nil {
		r.restOtherErrorStatus(err)
		return r
	}

	if r.stdErr, err = r.cmd.StderrPipe(); err != nil {
		r.restOtherErrorStatus(err)
		return r
	}

	if err = r.cmd.Start(); err != nil {
		r.restErrorStatus(ResultStatusRunError, err)
		return r
	}

	r.status = ResultStatusRunning

	go func() {
		defer close(r.runOk)

		if err = r.cmd.Wait(); err != nil {
			r.restErrorStatus(ResultStatusRunWaitError, err)
			return
		}
		r.err = nil
		r.status = ResultStatusSuccess
	}()

	return r
}

// restStatusAndCallFn 重置状态并调用函数
func (r *Result) restStatusAndCallFn(status ResultStatus, fn func()) {
	r.status = status
	if fn != nil {
		fn()
	}
}

// restOtherErrorStatus 重置错误状态
func (r *Result) restOtherErrorStatus(err error) {
	r.restErrorStatus(ResultStatusOtherError, err)
}

// restErrorStatus 重置错误状态
func (r *Result) restErrorStatus(status ResultStatus, err error) {
	r.restStatusAndCallFn(status, func() {
		r.err = err
	})
}

func (r *Result) statusCheckAndCallback(callback func() error, status ...ResultStatus) error {
	if r.TryLock() {
		defer r.Unlock()
	}

	if r.IsError() {
		return r.Error()
	}

	if !slices.Contains(status, r.status) {
		return ErrResultStatusNoMatch
	}

	return callback()
}

func (r *Result) Pid() int {
	return r.cmd.ProcessState.Pid()
}

func (r *Result) String() (string, error) {
	r.Lock()
	defer r.Unlock()
	<-r.runOk

	var result string

	if err := r.statusCheckAndCallback(func() error {
		result = r.err.Error()
		return nil
	}, ResultStatusOtherError, ResultStatusRunError); !errors.Is(err, ErrResultStatusNoMatch) {
		return result, err
	}

	if err := r.statusCheckAndCallback(func() error {
		buf := &strings.Builder{}
		buf.WriteString(r.err.Error())
		if output, err := io.ReadAll(r.stdOut); err == nil {
			buf.WriteString("\nOUTPUT:\n")
			buf.Write(output)
		}

		if errOutput, err := io.ReadAll(r.stdErr); err != nil {
			buf.WriteString("\nERROR_OUTPUT:\n")
			buf.Write(errOutput)
		}

		result = buf.String()
		return nil
	}, ResultStatusRunWaitError); !errors.Is(err, ErrResultStatusNoMatch) {
		return result, err
	}

	if err := r.statusCheckAndCallback(func() error {
		if buf, err := io.ReadAll(r.stdOut); err != nil {
			return err
		} else {
			result = string(buf)
		}
		return nil
	}, ResultStatusSuccess); !errors.Is(err, ErrResultStatusNoMatch) {
		return result, err
	}

	return "", ErrResultStatusNoMatch
}
