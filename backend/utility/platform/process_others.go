//go:build !windows

//go:dev !windows

package platform

import (
	"os/exec"
)

// 在非 Windows 平台不修改子程序設定
func HideWindow(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}

// chdman的名字
func CmdName() string {
	return "/opt/homebrew/bin/chdman"
}
