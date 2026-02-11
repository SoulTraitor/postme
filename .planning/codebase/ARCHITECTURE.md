# Architecture

**Analysis Date:** 2026-02-11

## Pattern Overview

**Overall:** Wails Desktop Application with Backend/Frontend Separation

PostMe is a desktop HTTP client built with Wails (Go + Vue 3). It follows a clean architecture pattern with distinct separation between backend (Go) and frontend (TypeScript/Vue), communicating via Wails' IPC mechanism.

**Key Characteristics:**
- Desktop-first architecture (Windows/Mac/Linux via Wails framework)
- Backend-driven data model with frontend as presentation layer
- Pinia stores for frontend state management with computed properties
- Repository pattern for data access layer
- Service layer for business logic orchestration
- Handler pattern for Wails IPC exposure

## Layers

**Frontend (Vue 3 + TypeScript):**
- Purpose: UI presentation, user interaction, state visualization
- Location: `frontend/src/`
- Contains: Vue components, Pinia stores, composables, TypeScript types, services/api layer
- Depends on: Wails runtime, backend handlers via IPC
- Used by: Desktop application window (Wails runtime)

**Backend API Layer (Go Handlers):**
- Purpose: Expose business logic to frontend via Wails IPC binding
- Location: `internal/handlers/`
- Contains: Handler structs with methods for frontend consumption (`RequestHandler`, `CollectionHandler`, `EnvironmentHandler`, `HistoryHandler`, `AppStateHandler`, `DialogHandler`)
- Depends on: Services layer for business logic, database layer
- Used by: Wails IPC binding mechanism (exposed in `main.go`)

**Backend Business Logic (Go Services):**
- Purpose: Implement domain logic and business rules
- Location: `internal/services/`
- Contains: Service classes orchestrating operations (`RequestService`, `CollectionService`, `EnvironmentService`, `HistoryService`, `HTTPClient`)
- Depends on: Repository layer for data access, models
- Used by: Handlers, other services

**Backend Data Access (Go Repositories):**
- Purpose: Abstract database operations
- Location: `internal/database/repository/`
- Contains: Repository classes with CRUD operations (`RequestRepository`, `CollectionRepository`, `FolderRepository`, `EnvironmentRepository`, `HistoryRepository`, `AppStateRepository`)
- Depends on: Database connection, models
- Used by: Services layer

**Backend Data Models:**
- Purpose: Define domain entities and their structure
- Location: `internal/models/`
- Contains: Struct definitions with JSON/DB tags (`Request`, `Collection`, `Folder`, `Environment`, `History`, `Response`, `AppState`, `KeyValue`)
- Depends on: Standard library
- Used by: All layers (handlers, services, repositories)

**Database Layer:**
- Purpose: SQLite connection management and schema migration
- Location: `internal/database/`
- Contains: Database initialization, connection pool management, schema migrations
- Depends on: sqlx (SQL toolkit), modernc.org/sqlite
- Used by: Repository layer

## Data Flow

**Request Execution Flow:**

1. Frontend: User enters HTTP method, URL, headers, params, body in RequestPanel
2. Frontend: Stores in active tab (Pinia `useTabsStore()`)
3. Frontend: User clicks "Send" → calls `api.executeRequest()`
4. Frontend API Service: Calls Wails-generated `RequestHandler.Execute()` IPC
5. Backend Handler: `RequestHandler.Execute()` receives params, calls `HTTPClient.Execute()`
6. Backend HTTPClient: Executes HTTP request with uTLS fingerprinting, returns response
7. Backend Handler: Saves to history via `HistoryService`, returns response
8. Frontend: Updates `responseStore` with response data, UI re-renders response panel

**Collection Hierarchy Flow:**

1. Backend: Collections contain Folders which contain Requests
2. Backend Service: `CollectionService.GetTree()` builds hierarchical structure
3. Frontend: Receives tree, stores in `useCollectionStore()`
4. Frontend: Sidebar renders tree recursively with expand/collapse
5. Frontend: Click request → loads into tab via `TabsStore.openRequest()`

**State Persistence Flow:**

1. Frontend: User makes changes to request/tabs/UI settings
2. Frontend: Computed properties (e.g., `isDirty`) auto-update
3. Frontend: Watchers trigger debounced saves (1000ms for app state, 500ms for tabs)
4. Frontend: Calls `api.updateAppState()` or `api.saveTabSession()`
5. Backend: Handlers write to database via repositories
6. Backend: On app startup, loads saved state to restore UI and tabs

**HTTP Client Flow:**

1. Receives `ExecuteRequest` with URL, method, headers, body, bodyType, timeout
2. Builds request: processes body (form-data, raw, etc.), applies headers
3. Handles compression: supports gzip, deflate, brotli response decompression
4. TLS Fingerprinting: Uses uTLS (Chrome_120 profile) to avoid detection
5. Proxy Support: Respects system proxy via `SetUseSystemProxy()`
6. Request Cancellation: Supports context-based cancellation per tab
7. Returns response with status, headers, body, duration

**State Management:**

- Frontend: Pinia stores hold all UI/request state (tabs, collections, responses, app settings)
- Persistence: Selected state synced to backend database (tab sessions, app state, window position)
- Backend: Database as source of truth for persistent data
- Initialization: App.vue `loadData()` loads everything from backend on startup with priority-based loading

## Key Abstractions

**Tab Concept:**
- Purpose: Multi-request workspace within single window
- Examples: `frontend/src/stores/tabs.ts`, `frontend/src/components/tabs/TabBar.vue`
- Pattern: Tab represents a request being edited. Stores request details (method, URL, headers, params, body) plus metadata (isDirty, originalState for change detection)

**Collection Tree:**
- Purpose: Hierarchical organization of requests
- Examples: `internal/services/collection_service.go`, `frontend/src/stores/collection.ts`
- Pattern: Collections → Folders → Requests structure, serialized as CollectionTree for transfer

**Environment Variables:**
- Purpose: Dynamic value substitution in requests
- Examples: `internal/models/environment.go`, `frontend/src/stores/environment.ts`
- Pattern: Global and environment-specific variables, substituted in URL/headers/body using `{{variableName}}` syntax

**Response Storage:**
- Purpose: Track response for currently active tab
- Examples: `frontend/src/stores/response.ts`, `internal/models/response.go`
- Pattern: Separate store tracks response per tab ID to enable multi-tab response viewing

**Request History:**
- Purpose: Track all executed requests for audit/repeat
- Examples: `internal/models/history.go`, `internal/services/history_service.go`
- Pattern: Every executed request auto-saved with response details, queryable by request ID

## Entry Points

**Backend Entry Point:**
- Location: `main.go`
- Triggers: Program execution
- Responsibilities: Initialize database, create handler instances, configure Wails window, bind handlers to IPC, manage app lifecycle (startup/shutdown)

**Frontend Entry Point:**
- Location: `frontend/src/main.ts`
- Triggers: Wails loads HTML/JS
- Responsibilities: Create Vue app instance, mount Pinia, mount to DOM

**Frontend App Root:**
- Location: `frontend/src/App.vue`
- Triggers: Vue application lifecycle
- Responsibilities: Load all data from backend, render main layout (sidebar, tab bar, request panel, response panel), manage keyboard shortcuts, handle window state changes, coordinate store updates

**Request Execution:**
- Location: `frontend/src/components/request/RequestPanel.vue` (UI trigger) → `internal/handlers/request_handler.go` → `internal/services/http_client.go`
- Triggers: User clicks Send button or Ctrl+Enter
- Responsibilities: Collect request params, execute via HTTP client, capture response, save to history, update response store

## Error Handling

**Strategy:** Multi-layer error propagation with fallback handling

**Patterns:**
- Backend handlers return `error` interface; frontend API layer checks and throws
- Frontend catches errors in async operations, displays toast notifications via `useToast()`
- Database operations retry once if locked (SQLite concurrency workaround) - see App.vue lines 365-377
- HTTP client errors (network, timeout, TLS) wrapped with context
- History saves continue even if response received fails (handler line 142)
- Tab session saves log individual failures but continue (App.vue lines 472-479)

**Examples:**
- `internal/handlers/request_handler.go` line 116-148: Execute returns error if HTTP client fails
- `frontend/src/App.vue` line 365-377: Retry logic for database locked errors
- `frontend/src/services/api.ts`: Type conversions can handle missing/null values gracefully

## Cross-Cutting Concerns

**Logging:** Console-based in frontend (performance timers visible in Dev Tools), silent in backend (standard Go error returns)

**Validation:**
- Frontend: Computed properties validate state (isDirty calculation tabs.ts:32-58)
- Backend: Repositories validate inputs before database operations
- HTTP Client: Validates URL format, timeout values

**Authentication:** Not implemented - app assumes localhost/internal HTTP traffic, no auth layer

**Context Management:**
- Backend: Request execution uses context.Context for timeout and cancellation
- Frontend: Wails runtime context used for window operations
- Tab-scoped cancellation tracked in RequestHandler via map of context.CancelFunc

**Concurrency:**
- Backend: sync.Mutex protects cancel funcs map in RequestHandler
- Frontend: Pinia stores are reactive, no explicit mutex needed
- Database: SQLite handles single-writer constraint; Wails framework serializes calls

## Data Serialization

**Backend → Frontend:**
- JSON marshaling via `json` struct tags in models
- Complex types like KeyValue arrays stored as JSON strings in DB, unmarshaled on retrieval
- Wails automatically converts Go structs to TypeScript types

**Frontend → Backend:**
- Wails-generated TypeScript bindings handle serialization
- Pinia stores track changes, send only changed data
- API service layer (api.ts) handles type conversions between Wails types and frontend types

