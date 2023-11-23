package gocommandinvoker

// RunnerOptions 命令行执行器选项
type RunnerOptions struct {
	// Args 命令行参数
	Args []string
	// Env 环境变量
	Env []string
	// Dir 工作目录
	Dir string
	// SysProcAttr 系统进程属性
	SysProcAttr interface{}
}
