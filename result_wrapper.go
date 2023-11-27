package gocommandinvoker

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"
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

// ResultStatus 结果状态
type ResultStatus uint8

const (
	// ResultStatusNoRun 未运行
	ResultStatusNoRun ResultStatus = iota
	// ResultStatusRunning 运行中
	ResultStatusRunning
	// ResultStatusRunError 运行错误
	ResultStatusRunError
	// ResultStatusRunWaitError 等待错误
	ResultStatusRunWaitError
	// ResultStatusOtherError 错误
	ResultStatusOtherError
)

// Result 命令调用结果
type Result struct {
	sync.Mutex
	cmd            *exec.Cmd
	err            error
	tempDir        string
	stdIn          *os.File
	stdOut         *os.File
	stdOutReadonly *os.File
	stdErr         *os.File
	stdErrReadonly *os.File
	status         ResultStatus
	runOk          chan struct{}
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
func (r *Result) run() {
	r.Lock()
	defer r.Unlock()

	if r.runOk != nil {
		return
	}

	r.runOk = make(chan struct{}, 1)

	var (
		err     error
		tempDir string
	)

	defer func() {
		if err != nil {
			close(r.runOk)
		}
	}()

	if r.status != ResultStatusNoRun {
		return
	}

	tempDir, err = os.MkdirTemp("", "gocommandinvoker*")
	if err != nil {
		r.restOtherErrorStatus(ErrCreateTempDirFailed)
		return
	}

	r.tempDir = tempDir

	if r.stdIn, _, err = r.openFileInTempDir("_stdin", false); err != nil {
		return
	}

	if r.stdOut, r.stdOutReadonly, err = r.openFileInTempDir("_stdout", true); err != nil {
		return
	}

	if r.stdErr, r.stdErrReadonly, err = r.openFileInTempDir("_stderr", true); err != nil {
		return
	}

	r.cmd.Stdin = r.stdIn
	r.cmd.Stdout = r.stdOut
	r.cmd.Stderr = r.stdErr

	if err = r.cmd.Run(); err != nil {
		r.restErrorStatus(ResultStatusRunError, err)
		return
	}

	r.status = ResultStatusRunning

	go func() {
		defer close(r.runOk)

		if err = r.cmd.Wait(); err != nil {
			r.restErrorStatus(ResultStatusRunWaitError, err)
		}
	}()
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

// openFileInTempDir 在临时目录中打开文件
func (r *Result) openFileInTempDir(filename string, createReadonlyStream bool) (*os.File, *os.File, error) {
	if r.tempDir == "" {
		return nil, nil, ErrTempDirNotExists
	}

	fPath := filepath.Join(r.tempDir, filename)
	f, err := os.OpenFile(fPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		r.restOtherErrorStatus(ErrCreateTempFileFailed)
		return nil, nil, err
	}

	var fR *os.File

	if createReadonlyStream {
		if fR, err = os.OpenFile(fPath, os.O_RDONLY, 0666); err != nil {
			_ = f.Close()
			r.restOtherErrorStatus(ErrCreateTempFileFailed)
			return nil, nil, err
		}
	}

	return f, fR, nil
}

func (r *Result) statusCheckAndCallback(callback func() error, status ...ResultStatus) error {
	r.Lock()
	defer r.Unlock()

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
		errStr := r.err.Error()
		stdOutStat, err := r.stdOutReadonly.Stat()
		stdErrStat, err := r.stdErrReadonly.Stat()
		return nil
	}, ResultStatusRunWaitError); !errors.Is(err, ErrResultStatusNoMatch) {
		return result, err
	}

	return "", ErrResultStatusNoMatch
}
