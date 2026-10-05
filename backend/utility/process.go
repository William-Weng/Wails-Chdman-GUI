package utility

import (
	"os/exec"
	"wails-chdman-gui/backend/utility/platform"
)

// 設定子程序的執行屬性，主要用於在 Windows 上隱藏子程序的控制台視窗
// 參數 cmd:
//   - 要執行的 *exec.Cmd 物件
//
// 回傳值:
//   - 設定後的 *exec.Cmd 物件（通常就是傳入的同一個指標）
func HideWindow(cmd *exec.Cmd) *exec.Cmd {
	return platform.HideWindow(cmd)
}
