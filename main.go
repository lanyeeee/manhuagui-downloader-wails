package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"manhuagui-downloader/backend/api"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	downloadApi := api.NewDownloadApi()
	exportApi := api.NewExportApi()
	pathApi := api.NewPathApi()
	settingsApi := api.NewSettingsApi()
	utilsApi := api.NewUtilsApi()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "manhuagui-downloader",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			downloadApi.Startup(ctx)
			exportApi.Startup(ctx)
			pathApi.Startup(ctx)
			settingsApi.Startup(ctx)
			utilsApi.Startup(ctx)
		},
		Bind: []interface{}{
			downloadApi,
			exportApi,
			pathApi,
			settingsApi,
			utilsApi,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
