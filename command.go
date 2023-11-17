package gocommandinvoker

// CommandInvoker 命令调用者
type CommandInvoker interface {
	// SettingInvokerPrefix 设置命令前缀
	SettingInvokerPrefix(prefixCmd []string) CommandInvoker
	// Invoke 调用, 传入需要执行的命令，返回程序的退出代码、错误
	Invoke(cmd string) (int, error)
	// InvokeWithOutput 调用, 传入需要执行的命令，返回程序的执行结果退出代码、错误
	InvokeWithOutput(cmd string) (*InvokeResult, int, error)
	InvokeRaw()
}
