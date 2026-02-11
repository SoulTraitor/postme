# Coding Conventions

**Analysis Date:** 2026-02-11

## Naming Patterns

**Files:**
- Go files use snake_case: `http_client.go`, `collection_repo.go`, `app_state.ts`
- Vue components use PascalCase: `KeyValueEditor.vue`, `BodyEditor.vue`, `ResponsePanel.vue`
- Utility/composable files use camelCase with `use` prefix: `useKeyboardActions.ts`, `useKeyValue.ts`, `useRequest.ts`
- Store files use camelCase: `appState.ts`, `collection.ts`, `environment.ts`
- Type definition files use index.ts in module: `src/types/index.ts`

**Functions/Methods:**
- Go: Use PascalCase for exported functions, camelCase for unexported: `NewRequestHandler()`, `Init()`, `GetByID()`
- TypeScript/Vue: Use camelCase for all functions: `convertKeyValue()`, `emitKeyboardAction()`, `onKeyboardAction()`
- Go receivers use short abbreviations: `func (r *CollectionRepository)`, `func (h *RequestHandler)`, `func (t *utlsTransport)`

**Variables:**
- Go: Use PascalCase for exported (package level), camelCase for unexported: `var DB *sqlx.DB`, `var defaultUserAgent`
- TypeScript/Vue: Use camelCase: `const cancelFuncs`, `const sidebarOpen`, `const effectiveTheme`
- Refs in Vue: Use camelCase: `const sidebarOpen = ref(true)`, `const localItems = ref([])`
- Boolean-prefixed names use `is`, `use`, `show`: `isModalOpen`, `useSystemProxy`, `showFileUpload`, `showTooltip`

**Types:**
- Go: Use PascalCase for all types: `type RequestHandler struct`, `type HTTPClient struct`, `type KeyValue struct`
- TypeScript interfaces: Use PascalCase: `interface Request`, `interface Collection`, `interface Environment`
- Type aliases/unions: Use camelCase: `type ActionCallback = () => void`

## Code Style

**Formatting:**
- Vue: Use 2-space indentation in templates and scripts (inferred from component structure)
- Go: Standard Go formatting (gofmt compatible)
- TypeScript: 2-space indentation (from Vite/Vue configuration)
- Line length: No explicit limit enforced; Vue templates use readable wrapping

**Vue Script Setup Style:**
- All Vue components use `<script setup lang="ts">` syntax exclusively
- Props defined with `defineProps<T>()` for type safety: `defineProps<{ params: KeyValue[] }>()`
- Emits defined with `defineEmits<T>()` for type safety: `defineEmits<{ 'update:params': [value: KeyValue[]] }>()`
- Imports organized: Vue/Pinia first, then components (@/ path), then types (@/types)

**Template Binding Patterns:**
- Use v-bind shorthand: `:class`, `:value`, `:key`
- Use v-on shorthand: `@click`, `@change`, `@input`, `@mouseenter`, `@mouseleave`
- Conditional CSS classes use ternary/array syntax consistently:
  ```html
  :class="[
    effectiveTheme === 'dark' ? 'dark-class' : 'light-class',
    otherCondition ? 'class1' : 'class2'
  ]"
  ```
- Use Tailwind utility classes for all styling (no inline styles except dynamic style bindings)

**CSS Classes:**
- Tailwind CSS exclusively for styling
- Responsive prefixes: `sm:`, `md:`, `lg:` for breakpoints
- Custom color tokens from tailwind.config.js: `bg-dark-surface`, `text-accent`, `border-dark-border`
- Flex/grid utilities: `flex`, `flex-col`, `items-center`, `justify-between`
- Spacing: `gap-2`, `mb-4`, `px-3`, `py-1.5`

## Import Organization

**Order:**
1. Standard library/framework imports (Go: stdlib; Vue: vue, pinia, etc.)
2. External dependencies
3. Internal packages/modules (using `@/` path alias)
4. Types (from `@/types`)

**Path Aliases:**
- `@/*` → `./src/*` (defined in tsconfig.json and vite.config.ts)
- `wailsjs/*` → `./wailsjs/*` (Wails auto-generated bindings)
- Example imports:
  ```typescript
  import { useAppStateStore } from '@/stores/appState'
  import type { AppState, KeyValue } from '@/types'
  import * as RequestHandler from '../../wailsjs/go/handlers/RequestHandler'
  ```

**Go Import Groups:**
- Standard library first
- External packages second
- Internal modules (`github.com/SoulTraitor/postme/internal/...`) third

## Error Handling

**Patterns:**
- Go: Explicit error checking with `if err != nil { return err }`
- Go: Errors propagated up to handlers; handlers decide on response
- Go: Handler methods return `error` as second return value following Go convention
- TypeScript: Use try-catch for async operations:
  ```typescript
  try {
    const api = await getApi()
    await api.setSidebarItemExpanded(itemType, itemId, newValue)
  } catch (error) {
    console.error('Failed to save sidebar state:', error)
  }
  ```
- Vue: Errors logged with `console.error()` with context prefix
- Silent failures acceptable with `console.error()` for non-critical UI state persistence

## Logging

**Framework:** `console.log()`, `console.error()`, `console.warn()` (browser native)

**Patterns:**
- Performance metrics prefixed: `console.log('[Performance] message')`
- Error context: `console.error('Failed to load background data:', err)`
- Tab operations: `console.log('[Tabs] message')`
- Critical issues: `console.error('[CRITICAL] message')`
- Info with measurements: `console.log('message in ${time.toFixed(2)}ms')`
- Database issues: `console.warn('Database locked, retrying in 500ms...')`

**When to Log:**
- Performance timing measurements on startup
- Failed data loads with error details
- Critical state persistence issues
- Database lock retries
- Non-critical operation completions (with counts)

## Comments

**When to Comment:**
- Document non-obvious business logic
- Explain workarounds or temporary solutions (rare in codebase)
- Document public API methods (Go: comment above type/function)

**JSDoc/TSDoc:**
- Go: Use comment blocks above exported types and functions:
  ```go
  // Collection represents a top-level container for requests
  type Collection struct {
  ```
- Go: Functions documented briefly above declaration:
  ```go
  // NewRequestHandler creates a new RequestHandler
  func NewRequestHandler() *RequestHandler {
  ```
- TypeScript: Minimal JSDoc; rely on TypeScript type inference
- Vue: Doc comments rare; prop/emit types self-document with `defineProps<T>()`

## Function Design

**Size:**
- Go: Methods tend toward single responsibility; average 40-60 lines including error handling
- TypeScript: Composables and service functions kept compact (20-50 lines)
- Vue components: Template/logic roughly balanced

**Parameters:**
- Go: Use structs for multiple parameters (e.g., `ExecuteRequestParams`)
- Go: Pointer receivers for methods that modify state
- TypeScript: Use objects/interfaces for multiple params in public functions
- Vue props: Always use `defineProps<T>()` for type-safe destructuring

**Return Values:**
- Go: Multiple returns using `(value, error)` pattern
- TypeScript: Async functions return `Promise<T>` or throw errors
- TypeScript services: Return union types for responses, empty promises for mutations
- Vue composables: Return object with functions/refs, never plain values

## Module Design

**Exports:**
- Vue components: Default export (SFC)
- Composables: Named exports for functions + composable factory: `export { onKeyboardAction, emitKeyboardAction, useKeyboardActions }`
- Services/stores: Export single instance (e.g., `export const api = { ... }`, `export default api`)
- Go packages: Export types and functions explicitly; main.go imports from internal/

**Barrel Files:**
- `src/types/index.ts` exports all type definitions
- No barrel files for components (import from full path)
- Stores imported directly: `import { useAppStateStore } from '@/stores/appState'`
- Composables imported directly from file: `import { useKeyboardActions } from '@/composables/useKeyboardActions'`

**Package Visibility:**
- Go internal/ packages: Only exported types/functions used by handlers
- Go handlers: Bind to Wails runtime, expose public methods
- Vue stores: One store per logical domain (appState, collection, environment, history, response, tabs)

---

*Convention analysis: 2026-02-11*
