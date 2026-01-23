package main

import (
	"context"
	"embed"

	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/handlers"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create handlers
	requestHandler := handlers.NewRequestHandler()
	collectionHandler := handlers.NewCollectionHandler()
	environmentHandler := handlers.NewEnvironmentHandler()
	historyHandler := handlers.NewHistoryHandler()
	appStateHandler := handlers.NewAppStateHandler()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "PostMe",
		Width:     1200,
		Height:    800,
		MinWidth:  800,
		MinHeight: 600,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 26, G: 26, B: 26, A: 1},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		OnStartup: func(ctx context.Context) {
			// Initialize database
			if err := database.Init(); err != nil {
				panic(err)
			}

			// Initialize handlers
			requestHandler.Init()
			collectionHandler.Init()
			environmentHandler.Init()
			historyHandler.Init()
			appStateHandler.Init()
		},
		OnShutdown: func(ctx context.Context) {
			database.Close()
		},
		Bind: []interface{}{
			requestHandler,
			collectionHandler,
			environmentHandler,
			historyHandler,
			appStateHandler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
