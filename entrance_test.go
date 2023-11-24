//go:build unix || (js && wasm) || plan9 || wasip1

package gocommandinvoker

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

// TestFindCommandPath 测试查找命令路径
func TestFindCommandPath(t *testing.T) {
	Convey("测试对外提供的入口函数", t, func() {

		testDataDir, e := filepath.Abs("./test_datas")
		if e != nil {
			panic(e)
		}

		Convey("测试查找存在的命令路径", func() {
			testCases := []struct {
				command string
				want    string
			}{
				{"ls", "/usr/bin/ls"},
				{"ls -la", "/usr/bin/ls"},
				{"pwd", "/usr/bin/pwd"},
				{"./test_datas/uid.sh", filepath.Join(testDataDir, "uid.sh")},
				{"./test_datas/gid.sh", filepath.Join(testDataDir, "gid.sh")},
				{"./test_datas/oid.sh", filepath.Join(testDataDir, "oid.sh")},
			}

			for _, testCase := range testCases {
				p, _, err := FindCommandPath(testCase.command)
				So(err, ShouldBeNil)
				So(p, ShouldEqual, testCase.want)
			}
		})

		Convey("测试查找不存在的命令路径", func() {
			testCases := []struct {
				command string
				err     error
			}{
				{"ls1", ErrNotFound},
				{"pwd1", ErrNotFound},
				{"./ls", ErrNotFound},
				{"../ls", ErrNotFound},
				{"~/ls", ErrNotFound},
				{"~/ls\\a -la", ErrNotFound},
				{filepath.Join(testDataDir, "no_execute.sh"), ErrNoExecutionPermissions},
			}

			for _, testCase := range testCases {
				_, _, err := FindCommandPath(testCase.command)
				So(errors.Is(err, testCase.err), ShouldBeTrue)
			}
		})

		Convey("测试查找空命令路径", func() {
			_, _, err := FindCommandPath("")
			So(errors.Is(err, ErrEmptyCommand), ShouldBeTrue)
		})
	})
}
