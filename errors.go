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
	// ErrCreateTempDirFailed 创建临时目录失败
	ErrCreateTempDirFailed = errors.New("create temp dir failed")
	// ErrCreateTempFileFailed 创建临时文件失败
	ErrCreateTempFileFailed = errors.New("create temp file failed")
	// ErrTempDirNotExists 临时目录不存在
	ErrTempDirNotExists = errors.New("temp dir not exists")
	// ErrResultStatusNoMatch 结果状态不匹配
	ErrResultStatusNoMatch = errors.New("result status no match")
)
