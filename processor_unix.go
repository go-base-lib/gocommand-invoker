//go:build unix || (js && wasm) || plan9 || wasip1

package gocommandinvoker

func init() {
	defaultCmdPrefix = "/usr/bin/sh -c"
	defaultRunner = New()
}
