package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"

	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/handlers"
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
	repo := repository.NewAppStateRepository(database.GetDB())
	savedState, _ := repo.Get()

	// Store context for shutdown
	var appCtx context.Context

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

			// Restore window position if saved and valid
			if savedState != nil && savedState.WindowX != nil && savedState.WindowY != nil {
				x, y := *savedState.WindowX, *savedState.WindowY
				// Only restore if position is reasonable (not negative or extremely large)
				// This prevents invisible windows from bad saved state
				if x >= -100 && x < 10000 && y >= -100 && y < 10000 {
					runtime.WindowSetPosition(ctx, x, y)
				}
			}
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			// Save window state before closing (window is still valid here)
			isMaximized := runtime.WindowIsMaximised(appCtx)
			isMinimized := runtime.WindowIsMinimised(appCtx)
			w, h := runtime.WindowGetSize(appCtx)
			x, y := runtime.WindowGetPosition(appCtx)
			
			// Get current state from database
			currentState, err := repo.Get()
			if err != nil {
				return false
			}
			
			// Update maximized state
			currentState.WindowMaximized = isMaximized
			
			// Save size/position only when normal (not maximized/minimized) and size is valid
			if !isMaximized && !isMinimized && w >= 800 && h >= 600 && x >= -100 && y >= -100 {
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
