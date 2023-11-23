package gocommandinvoker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var ()

func New() {

}

func NewWithPrefix(prefix string) {

}

var re = regexp.MustCompile(`\\.`)

var currentUser *user.User

func init() {
	var err error
	currentUser, err = user.Current()
	if err != nil {
		log.Fatalf("无法获取当前用户信息：%s", err)
	}
}

func hasExecutePermission(filename string) bool {
	if currentUser == nil {
		return false
	}

	info, err := os.Stat(filename)
	if err != nil {
		return false
	}

	uid, err := strconv.ParseUint(currentUser.Uid, 10, 32)
	if err != nil {
		return false
	}

	gid := info.Sys().(*syscall.Stat_t).Gid
	fileMode := info.Mode().Perm()

	if fileMode&os.ModePerm&0100 != 0 && uint32(uid) == info.Sys().(*syscall.Stat_t).Uid {
		// 文件所有者具有可执行权限
		return true
	}

	groups, err := os.Getgroups()
	if err != nil {
		return false
	}

	for _, groupID := range groups {
		if fileMode&os.ModePerm&0010 != 0 && uint32(groupID) == gid {
			// 文件所在组具有可执行权限
			return true
		}
	}

	// 其他用户具有可执行权限
	return fileMode&os.ModePerm&0001 != 0
}

// FindCommandPath 查找命令行路径
func FindCommandPath(cmd string) (string, error) {
	if currentUser == nil {
		return "", ErrCannotGetCurrentUser
	}
	// 去除首尾空格
	cmd = strings.TrimSpace(cmd)

	// 如果命令为空，返回错误
	if cmd == "" {
		return "", ErrEmptyCommand
	}

	// 使用正则表达式将命令分割成单词，以处理带有转义空格的情况
	cmd = re.ReplaceAllStringFunc(cmd, func(match string) string {
		var builder strings.Builder
		builder.Grow(len(match))
		for i := 0; i < len(match); i++ {
			if match[i] == '\\' {
				builder.WriteByte('\x00')
			} else {
				builder.WriteByte(match[i])
			}
		}
		return builder.String()
	})

	cmd = strings.ReplaceAll(strings.Fields(cmd)[0], "\x00", " ")

	// 如果 LookPath 找不到，尝试将命令作为相对路径解析
	if strings.HasPrefix(cmd, "./") || strings.HasPrefix(cmd, "../") || strings.HasPrefix(cmd, "~/") {
		// 如果命令以 "~/" 开头，替换为用户的主目录
		if strings.HasPrefix(cmd, "~/") {
			cmd = filepath.Join(currentUser.HomeDir, cmd[2:])
		}

		if p, err := filepath.Abs(cmd); err != nil {
			return "", err
		} else if hasExecutePermission(p) {
			return p, nil
		}
	}

	// 尝试使用 LookPath 查找命令
	if p, err := exec.LookPath(cmd); err == nil {
		return p, nil
	}

	return "", fmt.Errorf("%w: %s", ErrNotFound, cmd)
}
