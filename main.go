package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"

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

// getWebviewUserDataPath returns the path for WebView2 user data
func getWebviewUserDataPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(configDir, "postme")
}

func main() {
	// Create handlers
	requestHandler := handlers.NewRequestHandler()
	dialogHandler := handlers.NewDialogHandler()
	collectionHandler := handlers.NewCollectionHandler(dialogHandler)
	environmentHandler := handlers.NewEnvironmentHandler()
	historyHandler := handlers.NewHistoryHandler()
	appStateHandler := handlers.NewAppStateHandler()

	// Initialize database early to restore window state
	if err := database.Init(); err != nil {
		panic(err)
	}

	// Load saved app state for window restoration
	repo := repository.NewAppStateRepository(database.GetDB())
	savedState, _ := repo.Get()

	// Store context for shutdown
	var appCtx context.Context

	// Default window settings
	windowWidth := models.DefaultWindowWidth
	windowHeight := models.DefaultWindowHeight
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
			WebviewUserDataPath:  getWebviewUserDataPath(),
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		OnStartup: func(ctx context.Context) {
			appCtx = ctx
			// Database already initialized, just init handlers
			requestHandler.Init()
			collectionHandler.Init()
			environmentHandler.Init()
			historyHandler.Init()
			appStateHandler.Init()
			dialogHandler.SetContext(ctx)

			restoreSavedWindowBounds(ctx, savedState, windowWidth, windowHeight)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			windowCtx := appCtx
			if windowCtx == nil {
				windowCtx = ctx
			}

			// Save window state before closing (window is still valid here)
			isMaximized := runtime.WindowIsMaximised(windowCtx)
			isMinimized := runtime.WindowIsMinimised(windowCtx)
			w, h := runtime.WindowGetSize(windowCtx)
			x, y := runtime.WindowGetPosition(windowCtx)

			// Get current state from database
			currentState, err := repo.Get()
			if err != nil {
				return false
			}

			// Update maximized state
			currentState.WindowMaximized = isMaximized

			// Save size/position only when normal (not maximized/minimized) and size is valid
			if !isMaximized && !isMinimized && shouldSaveWindowBounds(windowCtx, x, y, w, h) {
				currentState.WindowWidth = w
				currentState.WindowHeight = h
				currentState.WindowX = &x
				currentState.WindowY = &y
			}

			// Save to database
			_ = repo.Update(currentState)

			return false
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
