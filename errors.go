package gocommandinvoker

import "errors"

var (
	// ErrEmptyCommand 空命令
	ErrEmptyCommand = errors.New("empty command")
	// ErrNotFound 未找到
	ErrNotFound = errors.New("not found")
	// ErrCannotGetCurrentUser 无法获取当前用户
	ErrCannotGetCurrentUser = errors.New("cannot get current user")
	// ErrNoExecutionPermissions 无执行权限
	ErrNoExecutionPermissions = errors.New("no execution permissions")
	// ErrResultStatusNoMatch 结果状态不匹配
	ErrResultStatusNoMatch = errors.New("result status no match")
	// ErrExecNoRun 执行未运行
	ErrExecNoRun = errors.New("exec no run")
)
