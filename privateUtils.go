package gocommandinvoker

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var (
	currentUser *user.User
	re          = regexp.MustCompile(`\\.`)
)

func init() {
	var err error
	currentUser, err = user.Current()
	if err != nil {
		log.Fatalf("无法获取当前用户信息：%s", err)
	}
}

// hasExecutePermission 判断文件是否具有可执行权限
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

// replaceEmptyStr 替换空字符串
func replaceEmptyStr(rawStr, emptyReplaceStr string) string {
	rawStr = strings.TrimSpace(rawStr)
	if len(rawStr) == 0 {
		return emptyReplaceStr
	}
	return rawStr
}

// isEmptyStr 判断字符串是否为空
func isEmptyStr(str string) bool {
	str = strings.TrimSpace(str)
	return len(str) == 0
}

func splitStrCmd(cmd string) []string {
	if isEmptyStr(cmd) {
		return nil
	}

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

	fields := strings.Fields(cmd)
	for i := range fields {
		fields[i] = strings.ReplaceAll(fields[i], "\x00", " ")
	}
	return fields
}

func handleSysPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if isEmptyStr(p) {
		return p, nil
	}

	if strings.HasPrefix(p, "~") {
		p = filepath.Join(currentUser.HomeDir, p[1:])
	}

	if strings.HasPrefix(p, "~/") {
		p = filepath.Join(currentUser.HomeDir, p[2:])
	}

	return filepath.Abs(p)
}
