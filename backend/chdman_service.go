package backend

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
	"wails-chdman-gui/backend/utility"
	util "wails-chdman-gui/backend/utility"
	"wails-chdman-gui/backend/utility/platform"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ChdmanService 負責處理 CHD 相關操作，例如解析檔案路徑、Extract CD，以及 Create CD
type ChdmanService struct {
	App    *application.App
	mutex  sync.Mutex
	cancel context.CancelFunc
}

var extractingRegexp = regexp.MustCompile(`Extracting,\s*([0-9]+(?:\.[0-9]+)?)%\s*complete`)
var creatingRegexp = regexp.MustCompile(`Compressing,\s*([0-9]+(?:\.[0-9]+)?)%\s*complete.*\(ratio=([0-9]+(?:\.[0-9]+)?)%\)`)

// ServiceShutdown 是 Wails v3 的 service lifecycle hook
//
// 當 Wails App 關閉時會自動呼叫，例如使用者按下視窗 X、Cmd+Q 或系統登出；這裡取消目前任務的 context，讓 exec.CommandContext 終止還在執行的 FFmpeg
func (service *ChdmanService) ServiceShutdown() error {
	return service.Cancel()
}

// ParseFilePath 解析輸入的檔案路徑，並將目錄、檔名與副檔名包裝成 ParseFilePathResult 回傳
func (service *ChdmanService) ParseFilePath(filepath string) ParseFilePathResult {
	dir, fileName, ext := util.ParseFilePath(filepath)
	return ParseFilePathResult{Dir: dir, FileName: fileName, Ext: ext}
}

// 使用 chdman 的 extractcd 指令，將 CHD 檔案解開成 CUE/BIN 檔案
func (service *ChdmanService) ExtractCD(inputPath string, useDVD bool) (string, error) {

	dir, fileName, _ := util.ParseFilePath(inputPath)
	_ = service.Cancel()

	if err := util.CheckInputFile(inputPath); err != nil {
		return "", fmt.Errorf("輸入檔案錯誤：%v", err)
	}

	ctx, cancel := util.CreateCancelContext(0)
	service.cancel = cancel
	defer cancel()

	args := utility.ExtractArgs(inputPath, dir, fileName, useDVD)

	fmt.Print(args)

	command := platform.CmdName()
	cmd, err := utility.CreateCommandContext(ctx, command, args)

	if err != nil {
		return "", fmt.Errorf("建立指令失敗：%v", err)
	}

	start := time.Now()

	err = utility.RunCommandStream(
		cmd,
		func(source string, text string) {

			if source != "stderr" {
				return
			}

			matches := extractingRegexp.FindStringSubmatch(text)
			if len(matches) < 2 {
				return
			}

			progress := matches[1] + " %"

			service.App.Event.Emit("extract-progress", map[string]string{
				"progress": progress,
			})
		},
	)

	if err != nil {
		return "", fmt.Errorf("轉檔失敗：%v", err)
	}

	message := fmt.Sprintf("耗時：%s", time.Since(start).Round(time.Millisecond))
	return message, nil
}

// CreateCD 使用 chdman 的 createcd 指令，將 ISO/CUE/BIN 檔案建立成 CHD 檔案
//
// 回傳值：
//   - string：執行成功後的訊息，例如處理耗時
//   - error：執行過程中發生的錯誤
func (service *ChdmanService) CreateCD(inputPath string, useDVD bool) (string, error) {

	dir, fileName, _ := util.ParseFilePath(inputPath)
	_ = service.Cancel()

	if err := util.CheckInputFile(inputPath); err != nil {
		return "", fmt.Errorf("輸入檔案錯誤：%v", err)
	}

	ctx, cancel := util.CreateCancelContext(0)
	service.cancel = cancel
	defer cancel()

	args := utility.CreateArgs(inputPath, dir, fileName, useDVD)
	command := platform.CmdName()
	cmd, err := utility.CreateCommandContext(ctx, command, args)

	if err != nil {
		return "", fmt.Errorf("建立指令失敗：%v", err)
	}

	start := time.Now()
	ratio := "0.0 %"

	err = utility.RunCommandStream(
		cmd,
		func(source string, text string) {

			if source != "stderr" {
				return
			}

			matches := creatingRegexp.FindStringSubmatch(text)
			if len(matches) < 3 {
				return
			}

			progress := matches[1] + " %"
			ratio = matches[2] + " %"

			service.App.Event.Emit("create-progress", map[string]string{
				"progress": progress,
				"ratio":    ratio,
			})
		},
	)

	if err != nil {
		return "", fmt.Errorf("轉檔失敗：%v", err)
	}

	message := fmt.Sprintf("耗時：%s, 壓縮率：%s", time.Since(start).Round(time.Millisecond), ratio)
	return message, nil
}

// Cancel 供 Wails 前端呼叫，主動停止目前正在執行的 chdman 任務
func (service *ChdmanService) Cancel() error {

	service.mutex.Lock()
	defer service.mutex.Unlock()

	if service.cancel != nil {
		service.cancel()
		return nil
	}

	return errors.New("no running compose task")
}

// 檢查指定的 CHD 檔是否為 DVD-ROM 型 CHD => chdman info -i <chdPath>
//
// 參數:
//   - chdPath: CHD 檔的完整路徑
//
// 回傳值：
//   - true：CHD 的 metadata tag 為 "DVD "，代表 DVD-ROM CHD
//   - false：不是 DVD 型 CHD，通常可視為 CD 型 CHD
//   - error：無法執行 chdman、讀取輸出，或 CHD 檔無效
func (service *ChdmanService) CheckDVD(chdPath string) (bool, error) {

	command := platform.CmdName()
	out, err := exec.Command(command, "info", "-i", chdPath).Output()

	if err != nil {
		return false, fmt.Errorf("無法讀取 CHD 資訊: %w", err)
	}

	return strings.Contains(string(out), "Tag='DVD '"), nil
}
