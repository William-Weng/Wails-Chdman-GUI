package utility

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// 將給定的檔案路徑拆解為三個部分：目錄路徑、檔名（不含副檔名）、副檔名
// 參數：
//   - path：檔案的完整路徑，例如 "/Users/ios/Desktop/古本屋.jpg" 或 "data/config.yaml"
//
// 回傳值：
//   - "/Users/ios/Desktop/古本屋.jpg" → ("/Users/ios/Desktop", "古本屋", ".jpg")
//   - "/Users/ios/Desktop/.gitignore" → ("/Users/ios/Desktop", ".gitignore", "")
func ParseFilePath(path string) (rootdir string, filename string, extension string) {

	rootdir = filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)

	if ext == base {
		extension = ""
		filename = base
	} else {
		extension = ext
		filename = base[:len(base)-len(ext)]
	}

	return rootdir, filename, extension
}

// 讀取指定路徑的檔案內容，並以位元組切片（[]byte）回傳
//
// 參數：
//   - path：要讀取的檔案完整路徑，例如 "/Users/ios/.config/myapp/config.json"
//
// 回傳值：
//   - []byte：檔案的原始位元組內容。若檔案為空則回傳長度為 0 的切片。
//   - error：若讀取失敗則回傳錯誤。常見錯誤包括：
//   - os.ErrNotExist：檔案不存在
//   - 權限不足、路徑為目錄、或其他 I/O 錯誤
func ReadFile(path string) ([]byte, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// CheckInputFile 檢查輸入路徑是否為可存取的一般檔案 (路徑必須是一般檔案，而不是資料夾、symbolic link、socket、named pipe 或裝置檔)
//
// 參數：
//   - inputPath：要檢查的輸入檔案完整路徑，例如 "/Users/ios/Desktop/input.mp4"
func CheckInputFile(inputPath string) error {

	if inputPath == "" {
		return errors.New("input path cannot be empty")
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("input file does not exist: %q", inputPath)
		}

		return fmt.Errorf("cannot access input file %q: %w", inputPath, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("input path is not a regular file: %q", inputPath)
	}

	return nil
}
