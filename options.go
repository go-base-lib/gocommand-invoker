package gocommandinvoker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// RunnerOptions 命令行执行器选项
type RunnerOptions struct {
	// Args 命令行参数
	Args []string
	// Env 环境变量
	Env []string
	// Dir 工作目录
	Dir string
	// SysProcAttr 系统进程属性
	SysProcAttr *syscall.SysProcAttr
	// Ctx 上下文
	Ctx context.Context
	// cmdStr 命令
	cmdStr string
	// prefix 命令前缀
	prefix string
}

// Command 命令
func (r *RunnerOptions) Command() string {
	return r.cmdStr
}

// fillDefault 填充默认值
func (r *RunnerOptions) fillDefault(cmd string) error {
	if isEmptyStr(cmd) {
		return ErrEmptyCommand
	}

	r.prefix = replaceEmptyStr(r.prefix, defaultCmdPrefix)

	p, args, err := FindCommandPath(r.prefix)
	if err != nil {
		return fmt.Errorf("%w: 前缀命令[%s]未能识别", err, r.prefix)
	}
	r.prefix = p

	runnerArgs := args

	p, _, err = FindCommandPath(cmd)
	if err != nil {
		return err
	}

	if len(r.Args) > 0 {
		cmd += " " + strings.Join(r.Args, " ")
	}

	r.Args = append(runnerArgs, cmd)

	r.cmdStr = p

	if r.Dir, err = handleSysPath(replaceEmptyStr(r.Dir, filepath.Dir(p))); err != nil {
		return err
	}

	r.Env = append(os.Environ(), r.Env...)
	return nil
}

func (r *RunnerOptions) generatorExecCmd() *exec.Cmd {
	var cmd *exec.Cmd
	if r.Ctx == nil {
		cmd = exec.Command(r.prefix, r.Args...)
	} else {
		cmd = exec.CommandContext(r.Ctx, r.prefix, r.Args...)
	}
	cmd.Env = r.Env
	cmd.Dir = r.Dir
	cmd.SysProcAttr = r.SysProcAttr
	return cmd
}
