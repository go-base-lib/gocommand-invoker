package gocommandinvoker

import "errors"

var (
	// ErrEmptyCommand 空命令
	ErrEmptyCommand = errors.New("empty command")
	// ErrNotFound 未找到
	ErrNotFound = errors.New("not found")
	// ErrCannotGetCurrentUser 无法获取当前用户
	ErrCannotGetCurrentUser = errors.New("cannot get current user")
)
