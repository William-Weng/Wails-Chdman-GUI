package main

import (
	"embed"
	"wails-chdman-gui/backend"

	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {}

func main() {

	app := application.New(application.Options{
		Name:        "Chdman-GUI",
		Description: "一個使用 Wails 3、Go、Svelte 與 Less 製作的CHD檔案轉換工具",
		Services: []application.Service{
			application.NewService(&backend.ChdmanService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "CHD檔案轉換工具",
		Width:  256,
		Height: 256,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
		},
		EnableFileDrop:             true,
		DisableResize:              true,
		DefaultContextMenuDisabled: true,
		ZoomControlEnabled:         false,
		DevToolsEnabled:            false,
		BackgroundColour:           application.NewRGB(6, 7, 15),
		URL:                        "/",
	}).OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {

		ctx := event.Context()
		files := ctx.DroppedFiles()

		if len(files) > 0 {

			emitKey := "image-file-dropped"
			emitData := map[string]any{
				"path": files[0],
			}

			application.Get().Event.Emit(emitKey, emitData)
		}
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
