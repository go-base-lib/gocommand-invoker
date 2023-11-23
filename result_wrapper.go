package gocommandinvoker

func newErrResult(err error) *Result {
	return &Result{err: err}
}

// Result 命令调用结果
type Result struct {
	err error
}

func (r *Result) IsError() bool {
	return r.err != nil
}

func (r *Result) Error() error {
	return r.err
}
