package gocommandinvoker

import "errors"

var (
	// ErrEmptyCommand 空命令
	ErrEmptyCommand = errors.New("empty command")
	// ErrNotFound 未找到
	ErrNotFound = errors.New("not found")
)
