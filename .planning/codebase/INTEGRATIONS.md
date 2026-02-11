# External Integrations

**Analysis Date:** 2026-02-11

## APIs & External Services

**HTTP Requests (Primary):**
- User-initiated HTTP requests to arbitrary URLs
  - No predefined external API integrations
  - Backend executes user-specified HTTP requests with full protocol support
  - HTTP/1.1, HTTP/2, and custom TLS fingerprinting via uTLS
  - Location: `internal/services/http_client.go`

**System Proxy:**
- Uses system-configured proxy settings
  - Windows proxy detection via `golang.org/x/sys/windows/registry`
  - Environment variable fallback (standard HTTP_PROXY, HTTPS_PROXY)
  - Configurable via UI: `useSystemProxy` in app state
  - User can toggle proxy usage in settings

**TLS Fingerprinting (uTLS):**
- Spoofs browser TLS fingerprints to bypass detection systems
  - Chrome 120 fingerprint used by default
  - Package: `github.com/refraction-networking/utls`
  - Allows requests to protected endpoints (e.g., cloudflare)
  - Location: `internal/services/http_client.go` lines 90-93

## Data Storage

**Databases:**
- SQLite 3 (Pure Go implementation)
  - Connection: `modernc.org/sqlite` driver
  - File location: `{executable-directory}/data/postme.db`
  - Foreign keys enforced: `?_pragma=foreign_keys(1)`
  - ORM/Client: sqlx (struct scanning and named parameters)
  - Connection: `internal/database/db.go`

**Database Repositories:**
- `internal/database/repository/` contains:
  - `request_repo.go` - HTTP request CRUD
  - `collection_repo.go` - Collection organization
  - `environment_repo.go` - Environment and variable storage
  - `history_repo.go` - Request history persistence
  - `folder_repo.go` - Folder organization within collections
  - `app_state_repo.go` - Application state (window position, theme, settings)

**File Storage:**
- Local filesystem only
  - Database file: `data/postme.db`
  - WebView2 cache: `{USER_CONFIG_DIR}/postme/` (Windows)
  - No cloud storage or remote file uploads

**Caching:**
- None detected
- All data persists in SQLite
- No Redis, Memcached, or similar services

## Authentication & Identity

**Auth Provider:**
- Custom/None - Application does not integrate external authentication
- No OAuth, SAML, or third-party identity providers
- No user login/authentication system
- Single-user desktop application

**API Authentication Support:**
- Users can manually specify authentication headers in requests:
  - Headers editor in UI: `frontend/src/components/request/HeadersEditor.vue`
  - Supports: Bearer tokens, Basic auth, custom headers
  - Environment variable substitution for secure credential storage
  - Location: `internal/models/environment.go` and `internal/services/request_service.go`

**Environment Variables:**
- Used for secret management (API keys, tokens)
  - `internal/models/environment.go` - Variable model with `secret` flag
  - Variables can be marked as "secret" to hide in UI
  - Substituted in requests before execution
  - Location: `frontend/src/components/EnvironmentSelector.vue`

## Monitoring & Observability

**Error Tracking:**
- None detected
- No integration with Sentry, DataDog, or similar

**Logging:**
- Console output only during development
- Errors logged to stdout/stderr
- No persistent error logging
- Window state restoration error handling in `main.go`

**Application Metrics:**
- Response metrics collected:
  - Status code: `internal/models/history.go`
  - Duration (milliseconds): `durationMs` field
  - Response size: `size` field in Response model
  - Stored in history table for later review

## CI/CD & Deployment

**Hosting:**
- Standalone desktop application
- Windows-only (uses WebView2, Windows registry detection)
- No server deployment
- Binary distribution model

**Build System:**
- Wails CLI (via `wails build`)
- Cross-compilation possible but Windows-primary
- Output: Single executable `postme.exe`

**CI Pipeline:**
- Not detected
- No GitHub Actions, GitLab CI, Jenkins configuration found
- Manual build and release process implied

**Build Configuration:**
- `wails.json` - Frontend build and dev commands
- Frontend: `npm install`, `npm run build`, `npm run dev`
- Backend: Standard Go build via Wails

## Environment Configuration

**Required env vars:**
- None for application startup
- Optional:
  - `HTTP_PROXY`, `HTTPS_PROXY` - System proxy configuration (standard)
  - System theme preference (detected via OS settings, not env var)

**Secrets Location:**
- Secrets stored in SQLite database (`data/postme.db`)
  - Environment variables marked with `secret: true` flag
  - Not encrypted at rest
  - User responsibility to protect database file
- Backend server credentials stored as part of Requests/Environments
  - No separate secrets vault

**Windows-Specific Configuration:**
- WebView2 user data path: `{USER_CONFIG_DIR}/postme/`
- Registry check for system proxy settings: HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings

## Webhooks & Callbacks

**Incoming:**
- None detected
- Application is not an HTTP server receiving webhooks

**Outgoing:**
- None pre-configured
- Users can send HTTP requests to webhook endpoints manually
- No automatic webhook triggering

## Request Execution Features

**Compression Support:**
- Brotli decompression: `github.com/andybalholm/brotli`
- GZIP decompression: `compress/gzip` (stdlib)
- DEFLATE decompression: `compress/flate` (stdlib)
- Auto-decompression of server responses

**HTTP Features:**
- Full HTTP/1.1 support
- HTTP/2 support with ALPN negotiation
- Custom headers support
- Query parameters
- Request body with multiple types:
  - none
  - form-data
  - x-www-form-urlencoded
  - raw (JSON, XML, plain text)
- Response body capture and display

**Network Features:**
- Timeout support (configurable per request)
- Request cancellation (via context)
- System proxy support
- TLS 1.2+ support
- Custom certificate verification possible

**Response Handling:**
- Status code capture
- Full response headers
- Response body (with syntax highlighting via CodeMirror)
- Duration measurement
- Response size calculation
- History persistence for later reference

## External Dependencies Communication

**Frontend to Backend:**
- Wails IPC (Inter-Process Communication)
- Generated bindings in `frontend/wailsjs/go/handlers/`
- Type-safe async function calls from Vue components
- Error handling with promise rejection

**Type Conversion Layer:**
- Frontend types in `frontend/src/types/index.ts`
- API service converts between Wails-generated types and frontend types
- Located in `frontend/src/services/api.ts`

---

*Integration audit: 2026-02-11*
