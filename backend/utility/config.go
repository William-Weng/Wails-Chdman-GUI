package utility

import (
	"errors"
	"os"
	"path/filepath"
)

// 根據給定的資料夾名稱與檔名，產生位於使用者設定目錄中的完整設定檔路徑，並確保該資料夾存在（若不存在則自動建立）
//
// 參數：
//   - foldername：應用或模組專屬的設定資料夾名稱，例如 "myapp" 或 "audio-tool"
//   - filename：設定檔的檔名（含副檔名），例如 "config.json" 或 "settings.yaml"
//
// 回傳值：
//   - string：完整的設定檔路徑，例如：
//   - macOS/Linux: ~/.config/myapp/config.json
//   - Windows: %APPDATA%\myapp\config.json
//   - error：若取得使用者設定目錄失敗或建立資料夾失敗時回傳錯誤
func ConfigFilePath(folderName string, fileName string) (string, error) {

	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	appDir := filepath.Join(configDir, folderName)

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, fileName), nil
}

// 建立應用程式的 JSON 設定檔；若檔案已存在，則不會覆寫既有內容，直接回傳該檔案路徑
//
// 參數：
//   - folderName：應用程式或模組專屬的設定資料夾名稱，例如 "Demo-App"
//   - fileName：要建立的 JSON 設定檔名稱（含副檔名），例如 "config.json"
//
// 回傳值：
//   - string：成功建立或已存在的 JSON 設定檔完整路徑，例如：
//   - macOS: /Users/ios/Library/Application Support/Demo-App/config.json
//   - Linux: /home/ios/.config/Demo-App/config.json
//   - Windows: C:\Users\ios\AppData\Roaming\Demo-App\config.json
//   - error：取得設定檔路徑、建立檔案或寫入預設 JSON 內容失敗時回傳的錯誤
func CreateConfigJSON(folderName string, fileName string) (string, error) {

	path, err := ConfigFilePath(folderName, fileName)

	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)

	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return path, nil
		}
		return "", err
	}

	defer file.Close()

	if _, err := file.WriteString("{}\n"); err != nil {
		return "", err
	}

	return path, nil
}
