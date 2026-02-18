# Collection Export/Import Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add the ability to export a single Collection to a `.postme` JSON file, and import a `.postme` file as a new Collection.

**Architecture:** Backend-centric approach using Go for serialization and file I/O, with Wails native dialogs for file selection. Export strips IDs/timestamps, import creates new database records. Frontend adds "Export" to Collection right-click menu and "Import" button to sidebar.

**Tech Stack:** Go (backend logic, file I/O), Wails v2 runtime dialogs, Vue 3 + TypeScript (frontend UI), SQLite via sqlx (database)

---

### Task 1: Create export data structures

**Files:**
- Create: `internal/models/export.go`

**Step 1: Create the export model file**

```go
package models

import "time"

// ExportFile is the top-level export file structure
type ExportFile struct {
	Version    int              `json:"version"`
	ExportedAt time.Time       `json:"exportedAt"`
	Collection ExportCollection `json:"collection"`
}

// ExportCollection represents a collection without IDs/timestamps
type ExportCollection struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Folders     []ExportFolder  `json:"folders"`
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
	Name      string     `json:"name"`
	Method    string     `json:"method"`
	URL       string     `json:"url"`
	Headers   []KeyValue `json:"headers"`
	Params    []KeyValue `json:"params"`
	Body      string     `json:"body"`
	BodyType  string     `json:"bodyType"`
	SortOrder int        `json:"sortOrder"`
}
```

**Step 2: Verify it compiles**

Run: `cd F:/go/workspace/postme && go build ./internal/models/...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/models/export.go
git commit -m "feat: add export data structures for collection import/export"
```

---

### Task 2: Add SaveFileDialog to DialogHandler

**Files:**
- Modify: `internal/handlers/dialog_handler.go`

**Step 1: Add SaveFileDialog method**

Add the following method after the existing `OpenFileDialog` method in `internal/handlers/dialog_handler.go`:

```go
// SaveFileDialog opens a native file save dialog
func (h *DialogHandler) SaveFileDialog(title string, defaultFilename string) (string, error) {
	return runtime.SaveFileDialog(h.ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PostMe Files (*.postme)",
				Pattern:     "*.postme",
			},
		},
	})
}
```

**Step 2: Verify it compiles**

Run: `cd F:/go/workspace/postme && go build ./internal/handlers/...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/handlers/dialog_handler.go
git commit -m "feat: add SaveFileDialog to DialogHandler"
```

---

### Task 3: Add export/import service methods to CollectionService

**Files:**
- Modify: `internal/services/collection_service.go`

**Step 1: Add GetCollectionTree method**

This method retrieves a single collection's full tree (collection + folders + requests). Add after the existing `GetTree` method:

```go
// GetCollectionTree retrieves a single collection's full tree
func (s *CollectionService) GetCollectionTree(id int64) (*CollectionTree, error) {
	collection, err := s.collectionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	folders, err := s.folderRepo.GetByCollectionID(id)
	if err != nil {
		return nil, err
	}

	allRequests, err := s.requestRepo.GetByCollectionID(id)
	if err != nil {
		return nil, err
	}

	// Separate requests by folder
	requestsByFolder := make(map[int64][]models.Request)
	var directRequests []models.Request
	for _, req := range allRequests {
		if req.FolderID != nil {
			requestsByFolder[*req.FolderID] = append(requestsByFolder[*req.FolderID], req)
		} else {
			directRequests = append(directRequests, req)
		}
	}

	var folderTrees []FolderTree
	for _, folder := range folders {
		folderTrees = append(folderTrees, FolderTree{
			Folder:   folder,
			Requests: requestsByFolder[folder.ID],
		})
	}

	return &CollectionTree{
		Collection: *collection,
		Folders:    folderTrees,
		Requests:   directRequests,
	}, nil
}
```

**Step 2: Add ExportCollection method**

Converts a `CollectionTree` into an `ExportFile`. Add after `GetCollectionTree`:

```go
// ExportCollection converts a collection tree to an export file structure
func (s *CollectionService) ExportCollection(id int64) (*models.ExportFile, error) {
	tree, err := s.GetCollectionTree(id)
	if err != nil {
		return nil, err
	}

	exportFile := &models.ExportFile{
		Version:    1,
		ExportedAt: time.Now(),
		Collection: models.ExportCollection{
			Name:        tree.Collection.Name,
			Description: tree.Collection.Description,
		},
	}

	// Convert folders
	for _, ft := range tree.Folders {
		exportFolder := models.ExportFolder{
			Name:      ft.Folder.Name,
			SortOrder: ft.Folder.SortOrder,
		}
		for _, req := range ft.Requests {
			exportFolder.Requests = append(exportFolder.Requests, convertToExportRequest(req))
		}
		exportFile.Collection.Folders = append(exportFile.Collection.Folders, exportFolder)
	}

	// Convert direct requests
	for _, req := range tree.Requests {
		exportFile.Collection.Requests = append(exportFile.Collection.Requests, convertToExportRequest(req))
	}

	return exportFile, nil
}

func convertToExportRequest(req models.Request) models.ExportRequest {
	return models.ExportRequest{
		Name:      req.Name,
		Method:    req.Method,
		URL:       req.URL,
		Headers:   req.Headers,
		Params:    req.Params,
		Body:      req.Body,
		BodyType:  req.BodyType,
		SortOrder: req.SortOrder,
	}
}
```

**Step 3: Add ImportCollection method**

Creates a new collection from an `ExportFile`. Add after `ExportCollection`:

```go
// ImportCollection creates a new collection from an export file
func (s *CollectionService) ImportCollection(data *models.ExportFile) (*models.Collection, error) {
	// Create collection
	collection := &models.Collection{
		Name:        data.Collection.Name,
		Description: data.Collection.Description,
	}
	if err := s.collectionRepo.Create(collection); err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	// Create folders and their requests
	for _, ef := range data.Collection.Folders {
		folder := &models.Folder{
			CollectionID: collection.ID,
			Name:         ef.Name,
			SortOrder:    ef.SortOrder,
		}
		if err := s.folderRepo.Create(folder); err != nil {
			return nil, fmt.Errorf("failed to create folder %q: %w", ef.Name, err)
		}

		for _, er := range ef.Requests {
			req := &models.Request{
				CollectionID: collection.ID,
				FolderID:     &folder.ID,
				Name:         er.Name,
				Method:       er.Method,
				URL:          er.URL,
				Headers:      er.Headers,
				Params:       er.Params,
				Body:         er.Body,
				BodyType:     er.BodyType,
				SortOrder:    er.SortOrder,
			}
			if err := s.requestRepo.Create(req); err != nil {
				return nil, fmt.Errorf("failed to create request %q: %w", er.Name, err)
			}
		}
	}

	// Create direct requests (not in any folder)
	for _, er := range data.Collection.Requests {
		req := &models.Request{
			CollectionID: collection.ID,
			Name:         er.Name,
			Method:       er.Method,
			URL:          er.URL,
			Headers:      er.Headers,
			Params:       er.Params,
			Body:         er.Body,
			BodyType:     er.BodyType,
			SortOrder:    er.SortOrder,
		}
		if err := s.requestRepo.Create(req); err != nil {
			return nil, fmt.Errorf("failed to create request %q: %w", er.Name, err)
		}
	}

	return collection, nil
}
```

**Step 4: Add the `fmt` and `time` imports**

Make sure the import block at the top of `collection_service.go` includes `"fmt"` and `"time"`:

```go
import (
	"fmt"
	"time"

	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)
```

**Step 5: Verify it compiles**

Run: `cd F:/go/workspace/postme && go build ./internal/services/...`
Expected: No errors

**Step 6: Commit**

```bash
git add internal/services/collection_service.go
git commit -m "feat: add export/import service methods to CollectionService"
```

---

### Task 4: Add ExportCollection and ImportCollection to CollectionHandler

**Files:**
- Modify: `internal/handlers/collection_handler.go`

**Step 1: Add `ctx` field and dialog handler reference**

The `CollectionHandler` needs access to the Wails context (for dialogs) and the dialog handler. Modify the struct and its initialization:

Change the struct definition at the top of the file:

```go
// CollectionHandler handles collection-related operations for the frontend
type CollectionHandler struct {
	service *services.CollectionService
	dialog  *DialogHandler
}

// NewCollectionHandler creates a new CollectionHandler
func NewCollectionHandler(dialog *DialogHandler) *CollectionHandler {
	return &CollectionHandler{dialog: dialog}
}
```

**Step 2: Add ExportCollection method**

Add after the existing `ReorderRequests` method:

```go
// ExportCollection exports a collection to a .postme file
func (h *CollectionHandler) ExportCollection(id int64) error {
	// Get export data
	exportData, err := h.service.ExportCollection(id)
	if err != nil {
		return err
	}

	// Open save dialog
	defaultFilename := exportData.Collection.Name + ".postme"
	filePath, err := h.dialog.SaveFileDialog("Export Collection", defaultFilename)
	if err != nil {
		return err
	}
	if filePath == "" {
		return nil // User cancelled
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
```

**Step 3: Add ImportCollection method**

Add after `ExportCollection`:

```go
// ImportCollection imports a collection from a .postme file
func (h *CollectionHandler) ImportCollection() (*models.Collection, error) {
	// Open file dialog
	filePath, err := h.dialog.OpenFileDialog("Import Collection")
	if err != nil {
		return nil, err
	}
	if filePath == "" {
		return nil, nil // User cancelled
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON
	var exportFile models.ExportFile
	if err := json.Unmarshal(data, &exportFile); err != nil {
		return nil, fmt.Errorf("invalid file format: %w", err)
	}

	// Validate version
	if exportFile.Version != 1 {
		return nil, fmt.Errorf("unsupported file version: %d", exportFile.Version)
	}

	// Import into database
	collection, err := h.service.ImportCollection(&exportFile)
	if err != nil {
		return nil, err
	}

	return collection, nil
}
```

**Step 4: Update imports**

Add `"encoding/json"`, `"fmt"`, and `"os"` to the import block:

```go
import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/SoulTraitor/postme/internal/services"
)
```

**Step 5: Update OpenFileDialog filter in dialog_handler.go**

Modify the existing `OpenFileDialog` to also accept `.postme` files. Update the `Filters` in `internal/handlers/dialog_handler.go`:

```go
func (h *DialogHandler) OpenFileDialog(title string) (string, error) {
	return runtime.OpenFileDialog(h.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PostMe Files (*.postme)",
				Pattern:     "*.postme",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
}
```

**Step 6: Update main.go to pass DialogHandler to CollectionHandler**

In `main.go`, change the `NewCollectionHandler()` call to pass the dialog handler:

Change line 34 from:
```go
collectionHandler := handlers.NewCollectionHandler()
```
To:
```go
collectionHandler := handlers.NewCollectionHandler(dialogHandler)
```

**Step 7: Verify it compiles**

Run: `cd F:/go/workspace/postme && go build ./...`
Expected: No errors

**Step 8: Commit**

```bash
git add internal/handlers/collection_handler.go internal/handlers/dialog_handler.go main.go
git commit -m "feat: add ExportCollection and ImportCollection handler methods"
```

---

### Task 5: Regenerate Wails bindings

**Files:**
- Auto-generated: `frontend/wailsjs/go/handlers/CollectionHandler.d.ts`
- Auto-generated: `frontend/wailsjs/go/handlers/CollectionHandler.js`
- Auto-generated: `frontend/wailsjs/go/handlers/DialogHandler.d.ts`
- Auto-generated: `frontend/wailsjs/go/handlers/DialogHandler.js`

**Step 1: Generate Wails bindings**

Run: `cd F:/go/workspace/postme && wails generate module`
Expected: Bindings regenerated in `frontend/wailsjs/`

**Step 2: Verify new bindings exist**

Check that the following functions are now available:
- `CollectionHandler.ExportCollection(arg1: number): Promise<void>`
- `CollectionHandler.ImportCollection(): Promise<models.Collection | null>`
- `DialogHandler.SaveFileDialog(arg1: string, arg2: string): Promise<string>`

Run: `cat frontend/wailsjs/go/handlers/CollectionHandler.d.ts`

**Step 3: Commit**

```bash
git add frontend/wailsjs/
git commit -m "chore: regenerate Wails bindings for export/import"
```

---

### Task 6: Add API wrapper functions in frontend

**Files:**
- Modify: `frontend/src/services/api.ts`

**Step 1: Add exportCollection and importCollection to the api object**

Add after the `reorderRequests` method (around line 252) in the `api` object:

```typescript
  async exportCollection(id: number): Promise<void> {
    await CollectionHandler.ExportCollection(id)
  },

  async importCollection(): Promise<Collection | null> {
    const result = await CollectionHandler.ImportCollection()
    if (!result) return null
    return convertCollection(result)
  },
```

**Step 2: Commit**

```bash
git add frontend/src/services/api.ts
git commit -m "feat: add exportCollection and importCollection API wrappers"
```

---

### Task 7: Add Export menu item and Import button to CollectionTree.vue

**Files:**
- Modify: `frontend/src/components/sidebar/CollectionTree.vue`

**Step 1: Add icon imports**

In the icon import block (lines 149-159), add `ArrowDownTrayIcon` and `ArrowUpTrayIcon`:

```typescript
import {
  MagnifyingGlassIcon,
  PlusIcon,
  FolderPlusIcon,
  ChevronRightIcon,
  FolderIcon,
  FolderOpenIcon,
  PencilIcon,
  TrashIcon,
  DocumentDuplicateIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
} from '@heroicons/vue/24/outline'
```

**Step 2: Add "Export" to collectionMenuItems**

In the `collectionMenuItems` computed (around line 203), add the Export item before the Delete item:

```typescript
const collectionMenuItems = computed<ContextMenuItem[]>(() => [
  {
    id: 'add-folder',
    label: 'Add Folder',
    icon: FolderPlusIcon,
    action: () => addFolderToCollection(),
  },
  {
    id: 'add-request',
    label: 'Add Request',
    icon: PlusIcon,
    action: () => addRequestToCollection(),
  },
  {
    id: 'rename',
    label: 'Rename',
    icon: PencilIcon,
    action: () => renameCollection(),
  },
  {
    id: 'export',
    label: 'Export',
    icon: ArrowDownTrayIcon,
    action: () => exportCollection(),
  },
  {
    id: 'delete',
    label: 'Delete',
    icon: TrashIcon,
    danger: true,
    action: () => deleteCollection(),
  },
])
```

**Step 3: Add Import button to the sidebar top action area**

In the template section (around line 19-34), add an import button after the existing "New Folder" button:

```html
      <button
        @click="importCollection"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400' : 'hover:bg-light-hover text-gray-500'"
        title="Import Collection"
      >
        <ArrowUpTrayIcon class="w-5 h-5" />
      </button>
```

**Step 4: Add exportCollection and importCollection functions**

Add these functions in the script section, near the other collection action functions (like `deleteCollection`, `renameCollection`):

```typescript
async function exportCollection() {
  if (!selectedCollection.value) return
  try {
    await api.exportCollection(selectedCollection.value.id)
  } catch (error) {
    console.error('Failed to export collection:', error)
  }
}

async function importCollection() {
  try {
    const collection = await api.importCollection()
    if (collection) {
      await collectionStore.loadTree()
    }
  } catch (error) {
    console.error('Failed to import collection:', error)
  }
}
```

**Step 5: Commit**

```bash
git add frontend/src/components/sidebar/CollectionTree.vue
git commit -m "feat: add Export menu item and Import button to sidebar"
```

---

### Task 8: Build and manual test

**Step 1: Build the application**

Run: `cd F:/go/workspace/postme && wails build`
Expected: Build succeeds without errors

**Step 2: Manual test checklist**

1. Launch the app
2. Create a test collection with a folder and some requests
3. Right-click the collection → verify "Export" appears in context menu
4. Click "Export" → verify save dialog opens with `.postme` filter
5. Save to a known location → verify the file is valid JSON
6. Click the Import button in sidebar → verify open dialog appears
7. Select the exported file → verify a new collection appears in sidebar
8. Verify the imported collection contains all folders and requests

**Step 3: Final commit (if any fixes needed)**

```bash
git add -A
git commit -m "feat: collection export/import complete"
```
