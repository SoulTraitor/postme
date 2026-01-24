package main

import (
	"context"
	"embed"

	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/handlers"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	dialogHandler := handlers.NewDialogHandler()

	// Initialize database early to restore window state
	if err := database.Init(); err != nil {
		panic(err)
	}

	// Load saved app state for window restoration
	var savedState *models.AppState
	repo := repository.NewAppStateRepository(database.GetDB())
	savedState, _ = repo.Get()

	// Default window settings
	windowWidth := 1200
	windowHeight := 800
	windowStartState := options.Normal

	// Apply saved window settings if available
	if savedState != nil {
		if savedState.WindowWidth > 0 {
			windowWidth = savedState.WindowWidth
		}
		if savedState.WindowHeight > 0 {
			windowHeight = savedState.WindowHeight
		}
		if savedState.WindowMaximized {
			windowStartState = options.Maximised
		}
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "PostMe",
		Width:            windowWidth,
		Height:           windowHeight,
		MinWidth:         800,
		MinHeight:        600,
		Frameless:        true,
		WindowStartState: windowStartState,
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
			// Database already initialized, just init handlers
			requestHandler.Init()
			collectionHandler.Init()
			environmentHandler.Init()
			historyHandler.Init()
			appStateHandler.Init()
			dialogHandler.SetContext(ctx)

			// Restore window position if saved
			if savedState != nil && savedState.WindowX != nil && savedState.WindowY != nil {
				runtime.WindowSetPosition(ctx, *savedState.WindowX, *savedState.WindowY)
			}
		},
		OnShutdown: func(ctx context.Context) {
			database.Close()
		},
		Bind: []any{
			requestHandler,
			collectionHandler,
			environmentHandler,
			historyHandler,
			appStateHandler,
			dialogHandler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
