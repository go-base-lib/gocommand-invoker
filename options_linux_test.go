//go:build linux

package gocommandinvoker

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

// TestRunnerOptions 测试运行参数选项
func TestRunnerOptions(t *testing.T) {
	Convey("测试运行参数选项", t, func() {
		testCases := []struct {
			command       string
			err           error
			rawOptions    *RunnerOptions
			targetOptions *RunnerOptions
		}{
			{"", ErrEmptyCommand, &RunnerOptions{}, &RunnerOptions{}},
			{"ls -la", nil, &RunnerOptions{}, &RunnerOptions{
				cmdStr: "/usr/bin/ls",
				Args:   []string{"-c", "-la"},
				Env:    os.Environ(),
				Dir:    "/usr/bin",
				prefix: "/usr/bin/sh",
			}},
			{"lss", ErrNotFound, &RunnerOptions{}, &RunnerOptions{
				prefix: "/usr/bin/sh",
			}},
			{"ls -la", nil, &RunnerOptions{
				Dir: "~",
			}, &RunnerOptions{
				cmdStr: "/usr/bin/ls",
				Env:    os.Environ(),
				prefix: "/usr/bin/sh",
				Dir:    "~",
				Args:   []string{"-c", "-la"},
			}},
		}
		for _, testCase := range testCases {
			opt := testCase.rawOptions
			err := opt.fillDefault(testCase.command)
			So(errors.Is(err, testCase.err), ShouldBeTrue)
			So(opt, ShouldResemble, testCase.targetOptions)
		}
	})
}
