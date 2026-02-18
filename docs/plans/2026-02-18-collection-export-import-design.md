# Collection Export/Import Design

## Overview

Add the ability to export a single Collection (including its folders and requests) to a `.postme` JSON file, and import a `.postme` file as a new Collection.

## Scope

- Export: single Collection at a time, via Collection right-click context menu
- Import: single `.postme` file at a time, via sidebar top action area
- Format: PostMe-only JSON format (no Postman compatibility)
- Conflict: always create new Collection on import (no merge/overwrite)

## Export File Format

File extension: `.postme`

```json
{
  "version": 1,
  "exportedAt": "2026-02-18T12:00:00Z",
  "collection": {
    "name": "My API",
    "description": "API collection description",
    "folders": [
      {
        "name": "Users",
        "sortOrder": 0,
        "requests": [
          {
            "name": "Get Users",
            "method": "GET",
            "url": "https://api.example.com/users",
            "headers": [{"key": "Authorization", "value": "Bearer {{token}}", "enabled": true}],
            "params": [],
            "body": "",
            "bodyType": "none",
            "sortOrder": 0
          }
        ]
      }
    ],
    "requests": [
      {
        "name": "Health Check",
        "method": "GET",
        "url": "https://api.example.com/health",
        "headers": [],
        "params": [],
        "body": "",
        "bodyType": "none",
        "sortOrder": 0
      }
    ]
  }
}
```

Design decisions:
- All IDs and timestamps are stripped on export (not portable)
- `version` field for future format evolution
- `sortOrder` preserved to maintain ordering
- `exportedAt` for informational purposes

## Export Data Structures (Go)

```go
// ExportFile is the top-level export structure
type ExportFile struct {
    Version    int              `json:"version"`
    ExportedAt time.Time        `json:"exportedAt"`
    Collection ExportCollection `json:"collection"`
}

// ExportCollection represents a collection without IDs/timestamps
type ExportCollection struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    Folders     []ExportFolder `json:"folders"`
    Requests    []ExportRequest `json:"requests"`
}

// ExportFolder represents a folder without IDs/timestamps
type ExportFolder struct {
    Name      string          `json:"name"`
    SortOrder int             `json:"sortOrder"`
    Requests  []ExportRequest `json:"requests"`
}

// ExportRequest represents a request without IDs/timestamps
type ExportRequest struct {
    Name      string           `json:"name"`
    Method    string           `json:"method"`
    URL       string           `json:"url"`
    Headers   []models.KeyValue `json:"headers"`
    Params    []models.KeyValue `json:"params"`
    Body      string           `json:"body"`
    BodyType  string           `json:"bodyType"`
    SortOrder int              `json:"sortOrder"`
}
```

## Backend Changes

### 1. DialogHandler — Add SaveFileDialog

Add `SaveFileDialog(title string, defaultFilename string)` to `internal/handlers/dialog_handler.go`.

Uses `runtime.SaveFileDialog` from Wails with `.postme` file filter.

### 2. CollectionHandler — Export and Import Methods

Add to `internal/handlers/collection_handler.go`:

#### ExportCollection(id int64) error

1. Query collection tree for the given ID (reuse existing `GetTree` logic, filtered to one collection)
2. Open `SaveFileDialog` with default filename `{collection_name}.postme`
3. Convert `CollectionTree` to `ExportFile` (strip IDs/timestamps)
4. Marshal to indented JSON
5. Write to selected file path

#### ImportCollection() (*models.Collection, error)

1. Open `OpenFileDialog` with `.postme` file filter
2. Read file contents
3. Unmarshal JSON into `ExportFile`
4. Validate version field
5. Create Collection in database
6. Create Folders in database (map old references to new IDs)
7. Create Requests in database (assign to correct collection/folder IDs)
8. Return the newly created Collection

### 3. CollectionService — Add Import Support

Add `ImportCollection(data ExportFile)` method that handles database operations within a transaction to ensure atomicity.

## Frontend Changes

### 1. Collection Right-Click Menu — Add "Export"

In `CollectionTree.vue`, add to `collectionMenuItems`:

```typescript
{
  id: 'export',
  label: 'Export',
  icon: ArrowDownTrayIcon,
  action: () => exportCollection(),
}
```

### 2. Sidebar Top Action Area — Add "Import" Button

In `CollectionTree.vue`, add an import button next to the existing "New Collection" and "New Folder" buttons:

```html
<button @click="importCollection" title="Import Collection">
  <ArrowUpTrayIcon class="w-5 h-5" />
</button>
```

### 3. API Service — Add Wrapper Functions

In `api.ts`, add:
- `exportCollection(id: number): Promise<void>`
- `importCollection(): Promise<Collection | null>`

### 4. Post-Import — Refresh Tree

After successful import, call `collectionStore.loadTree()` to refresh the sidebar.

## Data Flow

### Export

```
User right-clicks Collection → "Export"
→ Frontend: api.exportCollection(id)
→ Wails: CollectionHandler.ExportCollection(id)
→ Go: Query CollectionTree for this ID
→ Go: SaveFileDialog → user picks path
→ Go: Convert to ExportFile, marshal JSON, write file
→ Return success
```

### Import

```
User clicks "Import" button in sidebar
→ Frontend: api.importCollection()
→ Wails: CollectionHandler.ImportCollection()
→ Go: OpenFileDialog → user picks .postme file
→ Go: Read file, unmarshal JSON, validate
→ Go: Create Collection + Folders + Requests in transaction
→ Return new Collection
→ Frontend: collectionStore.loadTree()
```

## Error Handling

- Invalid/corrupted file: return error message to frontend, display in toast/alert
- Unsupported version: return descriptive error
- User cancels file dialog: return nil/empty, frontend does nothing
- Database error during import: transaction rollback, return error

## Files to Create/Modify

### New Files
- `internal/models/export.go` — Export data structures

### Modified Files
- `internal/handlers/dialog_handler.go` — Add `SaveFileDialog`
- `internal/handlers/collection_handler.go` — Add `ExportCollection`, `ImportCollection`
- `internal/services/collection_service.go` — Add export/import service methods
- `frontend/src/components/sidebar/CollectionTree.vue` — Add menu item + import button
- `frontend/src/services/api.ts` — Add API wrapper functions
