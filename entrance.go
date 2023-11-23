package gocommandinvoker

import (
	"fmt"
	"os/exec"
	"strings"
)

var (
	defaultCmdPrefix = ""
)

func New() *Runner {
	return NewWithPrefix(defaultCmdPrefix)
}

func NewWithPrefix(prefix string) *Runner {
	return &Runner{
		prefix: prefix,
	}
}

// FindCommandPath 查找命令行路径
func FindCommandPath(cmd string) (string, []string, error) {
	if currentUser == nil {
		return "", nil, ErrCannotGetCurrentUser
	}

	cmdAndArgs := splitStrCmd(cmd)
	if len(cmdAndArgs) == 0 {
		return "", nil, ErrEmptyCommand
	}

	cmd = cmdAndArgs[0]

	// 如果 LookPath 找不到，尝试将命令作为相对路径解析
	if strings.HasPrefix(cmd, "./") || strings.HasPrefix(cmd, "../") || strings.HasPrefix(cmd, "~/") {
		if p, err := handleSysPath(cmd); err != nil {
			return "", nil, err
		} else if hasExecutePermission(p) {
			return p, nil, nil
		}
	}

	// 尝试使用 LookPath 查找命令
	if p, err := exec.LookPath(cmd); err == nil {
		return p, cmdAndArgs[1:], nil
	}

	return "", nil, fmt.Errorf("%w: %s", ErrNotFound, cmd)
}
