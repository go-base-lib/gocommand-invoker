package gocommandinvoker

type Runner interface {
	Exec(cmd string) Processor
	ExecWithOptions(cmd string)
}

type Processor interface {
}
