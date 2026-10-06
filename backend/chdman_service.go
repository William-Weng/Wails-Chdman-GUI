package backend

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"
	"wails-chdman-gui/backend/utility"
	util "wails-chdman-gui/backend/utility"
	"wails-chdman-gui/backend/utility/platform"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ChdmanService 負責處理 CHD 相關操作，例如解析檔案路徑、Extract CD，以及 Create CD
type ChdmanService struct {
	App *application.App // App 是 Wails 應用程式的實例；可用來和 Wails 的事件系統或應用程式生命週期互動
}

var extractingRegexp = regexp.MustCompile(`Extracting,\s*([0-9]+(?:\.[0-9]+)?)%\s*complete`)
var creatingRegexp = regexp.MustCompile(`Compressing,\s*([0-9]+(?:\.[0-9]+)?)%\s*complete.*\(ratio=([0-9]+(?:\.[0-9]+)?)%\)`)

// ParseFilePath 解析輸入的檔案路徑，並將目錄、檔名與副檔名包裝成 ParseFilePathResult 回傳
func (service *ChdmanService) ParseFilePath(filepath string) ParseFilePathResult {
	dir, fileName, ext := util.ParseFilePath(filepath)
	return ParseFilePathResult{Dir: dir, FileName: fileName, Ext: ext}
}

// 使用 chdman 的 extractcd 指令，將 CHD 檔案解開成 CUE/BIN 檔案
func (service *ChdmanService) ExtractCD(inputPath string) (string, error) {

	dir, fileName, _ := util.ParseFilePath(inputPath)

	if err := util.CheckInputFile(inputPath); err != nil {
		return "", fmt.Errorf("輸入檔案錯誤：%v", err)
	}

	ctx, cancel := util.CreateCancelContext(0)
	defer cancel()

	command := platform.CmdName()
	cmd, err := utility.CreateCommandContext(
		ctx,
		command,
		[]string{
			"extractcd",
			"-f",
			"-i", inputPath,
			"-o", filepath.Join(dir, fileName+".cue"),
			"-ob", filepath.Join(dir, fileName+".bin"),
		},
	)

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
func (service *ChdmanService) CreateCD(inputPath string) (string, error) {

	dir, fileName, _ := util.ParseFilePath(inputPath)

	if err := util.CheckInputFile(inputPath); err != nil {
		return "", fmt.Errorf("輸入檔案錯誤：%v", err)
	}

	ctx, cancel := util.CreateCancelContext(0)
	defer cancel()

	command := platform.CmdName()
	cmd, err := utility.CreateCommandContext(
		ctx,
		command,
		[]string{
			"createcd",
			"-f",
			"-i", inputPath,
			"-o", filepath.Join(dir, fileName+".chd"),
		},
	)

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
