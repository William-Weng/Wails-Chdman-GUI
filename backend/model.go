package backend

// 回傳給前端的資料結構
//
// Wails 會將這個 Go struct 轉換成前端可以使用的 JavaScript 物件，因此比直接回傳多個字串更適合用於 Wails binding
type ParseFilePathResult struct {
	Dir      string `json:"dir"`      // 檔案所在的目錄
	FileName string `json:"fileName"` // 不含副檔名的檔案名稱
	Ext      string `json:"ext"`      // 檔案副檔名
}
