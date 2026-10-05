package backend

import (
	"fmt"
	"path/filepath"
	"time"
	"wails-chdman-gui/backend/utility"
	util "wails-chdman-gui/backend/utility"
	"wails-chdman-gui/backend/utility/platform"
)

// ChdmanService 負責處理 CHD 相關操作，例如解析檔案路徑、Extract CD，以及 Create CD
type ChdmanService struct{}

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
			// fmt.Printf("[%s] %s", source, text)
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

	err = utility.RunCommandStream(
		cmd,
		func(source string, text string) {
			// fmt.Printf("[%s] %s", source, text)
		},
	)

	if err != nil {
		return "", fmt.Errorf("轉檔失敗：%v", err)
	}

	message := fmt.Sprintf("耗時：%s", time.Since(start).Round(time.Millisecond))
	return message, nil
}
