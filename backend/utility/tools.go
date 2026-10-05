package utility

import (
	"context"
)

type Tools struct {
}

// MARK: - Lifecycle Hooks
// ServiceStartup 是 Wails v3 的服務啟動生命週期勾子
//   - Wails 啟動應用程式時會呼叫此方法；這裡會取得目前的 Wails App 實例，讓 FFmpegService 後續可以透過 app 發送事件給前端
func (tools *Tools) ServiceStartup(ctx context.Context) error {
	return nil
}

// ServiceShutdown 在 Wails App 關閉時停止仍在執行的 FFmpeg
func (tools *Tools) ServiceShutdown() error {
	return nil
}

// MARK: - Public Methods
func NewTools() *Tools {
	return &Tools{}
}
