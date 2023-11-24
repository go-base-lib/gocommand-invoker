package gocommandinvoker

type Runner struct {
	// prefix 前缀
	prefix string
}

func (r *Runner) Exec(cmd string) *Processor {
	return r.ExecWithOptions(cmd, nil)
}

func (r *Runner) ExecWithOptions(cmd string, opt *RunnerOptions) *Processor {
	if opt == nil {
		opt = &RunnerOptions{}
	}

	opt.prefix = r.prefix
	opt.cmdStr = cmd

	return &Processor{
		RunnerOptions: opt,
	}
}

type Processor struct {
	*RunnerOptions
}

func (p *Processor) Run() *Result {
	if p.RunnerOptions == nil {
		return newErrResult(ErrNotFound)
	}

	if err := p.RunnerOptions.fillDefault(p.cmdStr); err != nil {
		return newErrResult(err)
	}

	return newResult(p.RunnerOptions.generatorExecCmd())
}
