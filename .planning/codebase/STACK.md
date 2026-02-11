# Technology Stack

**Analysis Date:** 2026-02-11

## Languages

**Primary:**
- TypeScript 5.8.3 - Frontend UI and type safety for Vue components and API layer
- Go 1.25.6 - Backend application logic, request handling, and desktop integration via Wails

**Secondary:**
- JavaScript - PostCSS and Tailwind configuration files

## Runtime

**Environment:**
- Node.js - Frontend development and build (via Vite)
- Go 1.25.6 - Desktop application runtime (Wails framework)
- Wails v2.11.0 - Desktop framework bridging Go backend with web frontend

**Package Manager:**
- npm - JavaScript/TypeScript dependencies (frontend)
  - Lockfile: Present (package-lock.json)
- go modules - Go dependencies
  - Lockfile: Present (go.sum)

## Frameworks

**Core Desktop Framework:**
- Wails v2.11.0 - Creates desktop application with Go backend + Vue3 frontend bridge
  - Generates bindings in `frontend/wailsjs/` for TypeScript-Go communication
  - Config: `wails.json`

**Frontend:**
- Vue 3.5.13 - UI framework for reactive components
  - Component template syntax with TypeScript support
  - Single-file components (.vue files)
  - Location: `frontend/src/`

**State Management:**
- Pinia 2.3.1 - Vue state management store for application state persistence
  - Stores: `frontend/src/stores/` (appState.ts, collection.ts, environment.ts, history.ts, response.ts, tabs.ts)

**Code Editor (In-App):**
- CodeMirror 6.0.1 - Embedded code/text editor for request bodies and responses
  - Language support via `@codemirror/lang-*` packages
  - Supports: HTML, JSON, XML
  - Theme: One Dark theme (`@codemirror/theme-one-dark`)

**Build/Dev:**
- Vite 6.2.4 - Frontend build tool and dev server
  - Config: `frontend/vite.config.ts`
  - Significantly faster than Webpack
  - ES module based development

**CSS Framework:**
- Tailwind CSS 3.4.17 - Utility-first CSS framework
  - Config: `frontend/tailwind.config.js`
  - Custom colors for HTTP methods, status codes, and themes
  - Dark mode support via `darkMode: 'class'`

**CSS Processing:**
- PostCSS 8.5.3 - CSS transformation
  - Config: `frontend/postcss.config.js`
  - Used with Autoprefixer for browser vendor prefixes

**Testing:**
- Not detected (no test framework or test files found)

**HTTP Client (Backend):**
- Standard library `net/http` - HTTP request/response handling
- `golang.org/x/net/http2` - HTTP/2 protocol support
- `github.com/refraction-networking/utls` v1.8.2 - uTLS for TLS fingerprint spoofing (evades detection)

## Key Dependencies

**Critical:**

**Frontend:**
- `@headlessui/vue` 1.7.23 - Headless UI components library for modal, dropdown, and menu elements
- `@heroicons/vue` 2.2.0 - SVG icon library for UI controls
- `@fontsource/jetbrains-mono` 5.2.8 - JetBrains Mono font (monospace for code/terminal display)
- `codemirror` 6.0.1 - Code editor kernel
  - `@codemirror/commands` 6.8.0 - Editor commands
  - `@codemirror/view` 6.36.4 - Editor rendering and DOM management
  - `@codemirror/state` 6.5.2 - Editor state management
  - `@codemirror/language` 6.10.8 - Language support framework
  - `@codemirror/search` 6.5.8 - Search functionality
  - `@codemirror/lang-html` 6.4.11 - HTML syntax highlighting
  - `@codemirror/lang-json` 6.0.1 - JSON syntax highlighting
  - `@codemirror/lang-xml` 6.1.0 - XML syntax highlighting
  - `@codemirror/theme-one-dark` 6.1.3 - Dark theme

**Backend:**
- `github.com/jmoiron/sqlx` v1.4.0 - SQL database wrapper with named query support
- `modernc.org/sqlite` v1.44.2 - Pure Go SQLite driver (in-process database)
- `github.com/andybalholm/brotli` v1.0.6 - Brotli compression support for HTTP responses
- `github.com/refraction-networking/utls` v1.8.2 - uTLS library for TLS fingerprinting
- `github.com/wailsapp/wails/v2` v2.11.0 - Wails framework
- `golang.org/x/net` v0.38.0 - Extended networking libraries
- `golang.org/x/sys` v0.37.0 - System-level networking and platform-specific code
- `golang.org/x/crypto` v0.36.0 - Cryptographic functions

**Infrastructure (Wails Dependencies - Indirect):**
- `github.com/wailsapp/go-webview2` v1.0.22 - WebView2 for Windows native web rendering
- `github.com/labstack/echo/v4` v4.13.3 - Web framework (used by Wails internally)
- `github.com/gorilla/websocket` v1.5.3 - WebSocket support for Wails IPC
- `github.com/google/uuid` v1.6.0 - UUID generation for unique identifiers
- `github.com/samber/lo` v1.49.1 - Go utility library

## Configuration Files

**TypeScript:**
- `frontend/tsconfig.json` - Strict mode enabled, module: ESNext
  - Path aliases: `@/*` → `./src/*`, `wailsjs/*` → `./wailsjs/*`
  - Node resolution for dependencies

**Build:**
- `frontend/vite.config.ts` - Vite configuration with Vue plugin, path aliases
- `frontend/tailwind.config.js` - Tailwind CSS configuration with custom color scheme
- `frontend/postcss.config.js` - PostCSS with Tailwind and Autoprefixer plugins

**Wails:**
- `wails.json` - Wails project configuration
  - App name: "postme"
  - Frontend build: `npm run build`
  - Frontend dev server: `npm run dev`

**Frontend Development:**
- `.gitignore` - Excludes node_modules, build artifacts, IDE files

## Environment Configuration

**Environment Variables:**
- No `.env` file detected
- No environment-based configuration variables found
- All configuration persists in SQLite database in `data/postme.db`

**Database Location:**
- Path: `{executable-directory}/data/postme.db`
- Fallback: `./data/` if executable path cannot be determined
- Created on first run via `database.Init()`

**WebView2 User Data:**
- Windows-specific: `{USER_CONFIG_DIR}/postme/`
- Used for WebView2 cache and user data storage

## Platform Requirements

**Development:**
- Go 1.25.6+
- Node.js (LTS recommended)
- npm 8+
- TypeScript 5.8.3+
- Windows SDK (for Windows build)
- Visual Studio Build Tools (for Windows compilation)

**Runtime:**
- Windows 10+ with WebView2 runtime
  - WebView2 is bundled with modern Windows installations
  - Fallback detection in code to handle user data path

## Type Safety

**TypeScript Compilation:**
- Strict mode enabled in tsconfig
- Type checking on build: `vue-tsc --noEmit`
- Wails auto-generates Go handler bindings in `frontend/wailsjs/go/`
- Generated types ensure type-safe backend communication

## Asset Serving

**Frontend:**
- Built files served from embedded FS (production)
  - `//go:embed all:frontend/dist` in `main.go`
- Dev server during `npm run dev` for live reload

---

*Stack analysis: 2026-02-11*
