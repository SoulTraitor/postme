# Codebase Structure

**Analysis Date:** 2026-02-11

## Directory Layout

```
postme/
├── frontend/                   # Vue 3 + TypeScript frontend application
│   ├── src/
│   │   ├── main.ts            # Entry point: creates Vue app + Pinia
│   │   ├── App.vue            # Root component: layout, lifecycle, data loading
│   │   ├── assets/            # Static assets and styles
│   │   │   └── styles/main.css
│   │   ├── components/        # Vue components organized by feature
│   │   │   ├── common/        # Reusable UI components
│   │   │   ├── request/       # Request editing components
│   │   │   ├── response/      # Response viewing components
│   │   │   ├── sidebar/       # Sidebar navigation
│   │   │   ├── tabs/          # Tab management UI
│   │   │   ├── modals/        # Dialog/modal components
│   │   │   ├── TitleBar.vue   # Custom window title bar
│   │   │   └── EnvironmentSelector.vue
│   │   ├── stores/            # Pinia state management
│   │   │   ├── appState.ts    # Global UI settings
│   │   │   ├── tabs.ts        # Open requests/tabs
│   │   │   ├── collection.ts  # Request hierarchy
│   │   │   ├── response.ts    # Response per tab
│   │   │   ├── environment.ts # Environment variables
│   │   │   └── history.ts     # Request history
│   │   ├── services/
│   │   │   └── api.ts         # Wails IPC API bindings + type conversion
│   │   ├── composables/       # Vue composition functions
│   │   │   ├── useRequest.ts  # Request building logic
│   │   │   ├── useKeyValue.ts # Key-value editor logic
│   │   │   ├── useKeyboardActions.ts # Keyboard shortcut emitter
│   │   │   └── useModal.ts    # Modal state management
│   │   └── types/
│   │       └── index.ts       # TypeScript type definitions
│   ├── wailsjs/               # Auto-generated Wails bindings (DO NOT EDIT)
│   │   ├── go/handlers/       # Generated Go handler stubs
│   │   └── runtime/           # Wails runtime API
│   ├── package.json           # Frontend npm dependencies
│   ├── vite.config.ts         # Vite build configuration
│   └── tsconfig.json          # TypeScript configuration
│
├── internal/                  # Backend Go code (not exposed to frontend directly)
│   ├── handlers/              # Wails IPC handler methods (bridge to frontend)
│   │   ├── request_handler.go
│   │   ├── collection_handler.go
│   │   ├── environment_handler.go
│   │   ├── history_handler.go
│   │   ├── app_state_handler.go
│   │   └── dialog_handler.go
│   │
│   ├── services/              # Business logic layer
│   │   ├── request_service.go
│   │   ├── collection_service.go
│   │   ├── environment_service.go
│   │   ├── history_service.go
│   │   └── http_client.go     # HTTP execution + uTLS
│   │
│   ├── database/              # Database initialization + migrations
│   │   ├── db.go              # SQLite connection, schema setup
│   │   ├── migrations.go      # Schema definitions (tables, indexes)
│   │   └── repository/        # Data access layer
│   │       ├── request_repo.go
│   │       ├── collection_repo.go
│   │       ├── folder_repo.go
│   │       ├── environment_repo.go
│   │       ├── history_repo.go
│   │       └── app_state_repo.go
│   │
│   └── models/                # Domain entity definitions
│       ├── request.go         # HTTP request definition
│       ├── collection.go      # Collection container
│       ├── folder.go          # Folder within collection
│       ├── environment.go     # Environment + variables
│       ├── history.go         # Executed request history
│       ├── response.go        # HTTP response
│       ├── app_state.go       # Persistent UI state
│       └── constants.go       # Enum-like constants
│
├── main.go                    # Application entry point: Wails setup + lifecycle
├── go.mod                     # Go module definition
├── go.sum                     # Go dependency lock file
├── wails.json                 # Wails configuration (frontend build, output)
│
├── build/                     # Build artifacts
│   ├── bin/                   # Compiled binaries
│   └── windows/               # Windows build resources
│
└── frontend/                  # Built/dist directory (not in git)
    └── dist/                  # Vite output (embedded in binary via //go:embed)
```

## Directory Purposes

**`frontend/src/`:**
- Purpose: All Vue 3 + TypeScript frontend source code
- Contains: Components, stores, services, types, assets
- Key files: `App.vue` (root), `main.ts` (entry), `services/api.ts` (backend bridge)

**`frontend/src/components/`:**
- Purpose: Vue component organization by UI domain
- Contains:
  - `common/` → Reusable widgets (KeyValueEditor, ToastContainer, ContextMenu)
  - `request/` → Request editing (BodyEditor, MethodSelect, UrlInput, ParamsEditor, HeadersEditor, RequestPanel)
  - `response/` → Response viewing
  - `sidebar/` → Collection tree navigation
  - `tabs/` → Tab bar and tab management
  - `modals/` → Dialog components

**`frontend/src/stores/`:**
- Purpose: Pinia reactive state management
- Contains: Centralized app state (tabs, collections, responses, environments, UI settings)
- Key pattern: Each store = domain (e.g., `tabs.ts` manages all tab state)

**`frontend/src/composables/`:**
- Purpose: Reusable Vue 3 composition functions
- Contains: `useRequest()` for request building, `useKeyValue()` for editor state, `useModal()` for dialogs

**`internal/handlers/`:**
- Purpose: Wails IPC method exposure - each handler is exposed as a JS object to frontend
- Pattern: Handler struct with methods → Wails auto-generates TypeScript stubs
- Binding: See `main.go` line 141-148: `Bind: []any{ requestHandler, collectionHandler, ... }`

**`internal/services/`:**
- Purpose: Business logic orchestration
- Pattern: Service receives repositories, implements use cases (e.g., CollectionService coordinates collection + folder + request repositories)
- Naming: `*Service` struct with public methods

**`internal/database/repository/`:**
- Purpose: Data access abstraction (CRUD per entity)
- Pattern: One repo per domain entity, methods marshal/unmarshal JSON for complex types
- Key detail: JSON fields stored as strings in SQLite (e.g., `headers` → `headersJSON` with JSON marshal/unmarshal)

**`internal/models/`:**
- Purpose: Domain entity definitions
- Pattern: Structs with `json` tags (for API/IPC) and `db` tags (for SQLite)
- Details: KeyValue is reusable type for headers/params/variables

## Key File Locations

**Entry Points:**
- `main.go`: Application bootstrap, Wails window setup, handler binding, lifecycle (startup/shutdown)
- `frontend/src/main.ts`: Vue app creation, Pinia initialization
- `frontend/src/App.vue`: Root component, data loading, layout rendering, global listeners

**Configuration:**
- `go.mod`: Go dependencies
- `go.sum`: Go dependency lock
- `wails.json`: Wails configuration (frontend dev/build commands)
- `frontend/package.json`: Node/npm dependencies
- `frontend/tsconfig.json`: TypeScript compiler options
- `frontend/vite.config.ts`: Vite bundler configuration (if exists)

**Core Logic:**
- `internal/services/http_client.go`: HTTP execution with uTLS, compression, proxy support
- `internal/services/collection_service.go`: Collection hierarchy logic (GetTree)
- `internal/database/migrations.go`: Schema definitions (tables, relationships)

**Database:**
- `internal/database/db.go`: Connection initialization, GetDB() singleton
- `internal/database/repository/*_repo.go`: CRUD operations

**Testing:**
- `frontend/src/types/index.ts`: Type definitions used by tests/components
- No test files currently in codebase

## Naming Conventions

**Files:**
- Go: `snake_case` for filenames (e.g., `app_state_handler.go`, `request_repo.go`)
- TypeScript/Vue: `snake_case` for utility files, `PascalCase` for components (e.g., `App.vue`, `RequestPanel.vue`, `useRequest.ts`)
- Stores: `*Store.ts` but exported as `use*Store()` (e.g., `tabs.ts` exports `useTabsStore()`)

**Directories:**
- Go: `snake_case` (e.g., `database`, `repository`, `handlers`)
- Frontend: `kebab-case` or descriptive (e.g., `components/request/`, `stores/`, `composables/`)

**Functions/Methods:**
- Go: `PascalCase` (public), `camelCase` (private) - e.g., `Create()`, `GetByID()`, `getDataDir()`
- TypeScript: `camelCase` (e.g., `useRequest()`, `buildUrl()`, `executeRequest()`)

**Types/Structs:**
- Go: `PascalCase` (e.g., `Request`, `Collection`, `RequestHandler`)
- TypeScript: `PascalCase` for interfaces/types (e.g., `Tab`, `Request`, `Collection`)

**Constants:**
- Go: `PascalCase` or `UPPER_SNAKE_CASE` (e.g., `defaultUserAgent`, `methodsWithoutBody`)
- TypeScript: Inline literals or `UPPER_SNAKE_CASE` for important values

**Variables:**
- Go: `camelCase` for receiver variables (e.g., `h`, `s`, `r` for handler/service/repo)
- TypeScript: `camelCase` for all (e.g., `activeTab`, `isDirty`, `isLoading`)

## Where to Add New Code

**New Feature (e.g., new request property):**
1. **Model**: Add field to `internal/models/request.go` with `json` + `db` tags
2. **Database**: Update `internal/database/migrations.go` to alter table
3. **Repository**: Update `internal/database/repository/request_repo.go` Create/Update/GetByID methods
4. **Service**: No change needed if repository handles it
5. **Handler**: No change needed if service exposes it
6. **Frontend Type**: Update `frontend/src/types/index.ts` Request interface
7. **Frontend Component**: Update editor component (e.g., `frontend/src/components/request/RequestPanel.vue`)
8. **Frontend Store**: Update `frontend/src/stores/tabs.ts` if it affects tab state

**New Component/Module:**
- Implementation: `frontend/src/components/{feature}/{ComponentName}.vue` or `frontend/src/{category}/{moduleName}.ts`
- If uses Pinia: Create `frontend/src/stores/{moduleName}.ts` with `defineStore()`
- If shared logic: Create `frontend/src/composables/use{ModuleName}.ts`

**New Backend Handler:**
1. Create `internal/handlers/{feature}_handler.go` with struct and Init/methods
2. Add service and repositories as needed
3. Bind in `main.go` Bind array
4. Wails auto-generates TypeScript stubs in `frontend/wailsjs/go/handlers/`
5. Create API wrapper in `frontend/src/services/api.ts`

**Utilities:**
- Shared frontend helpers: `frontend/src/services/` (utilities) or `frontend/src/composables/` (Vue hooks)
- Shared backend helpers: `internal/services/` (business logic) or separate utility file in relevant package

## Special Directories

**`frontend/wailsjs/`:**
- Purpose: Auto-generated Wails bindings (DO NOT MANUALLY EDIT)
- Generated: On `wails build` or when frontend hot-reloads
- Committed: Yes (committed so developers don't need to regenerate)
- Content: Go struct/interface stubs in TypeScript, Wails runtime API

**`build/`:**
- Purpose: Build artifacts
- Generated: Yes, on `wails build`
- Committed: Partially (windows resources committed, binaries not)
- Contains: Compiled binary, Windows resources (icons, manifests)

**`frontend/dist/`:**
- Purpose: Built frontend (Vite output)
- Generated: Yes, on `npm run build`
- Committed: No (ignored in .gitignore)
- Embedded: Via `//go:embed all:frontend/dist` in main.go

**`frontend/node_modules/`:**
- Purpose: npm dependencies
- Generated: Yes, on `npm install`
- Committed: No (ignored in .gitignore)
- Usage: Frontend build and dev

**`frontend/public/`:**
- Purpose: Static files (index.html, favicon)
- Generated: No
- Committed: Yes
- Content: Entry HTML, app resources

## File Organization Rules

**Keep related code close:**
- Request handler + request service + request repo should be in same vertical layer
- Request-related components (BodyEditor, ParamsEditor) co-located in `components/request/`

**Avoid deep nesting:**
- Frontend: Max 3 levels (e.g., `src/components/request/BodyEditor.vue`)
- Backend: Max 2 levels (e.g., `internal/database/repository/request_repo.go`)

**Separate concerns by layer:**
- UI logic → Vue components
- State management → Pinia stores
- Backend communication → `services/api.ts`
- Domain logic → handlers/services
- Data access → repositories

**Naming matches structure:**
- Import paths should read naturally: `from '@/components/request/RequestPanel.vue'` → path matches logic

