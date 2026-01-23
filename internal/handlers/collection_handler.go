package handlers

import (
	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/SoulTraitor/postme/internal/services"
)

// CollectionHandler handles collection-related operations for the frontend
type CollectionHandler struct {
	service *services.CollectionService
}

// NewCollectionHandler creates a new CollectionHandler
func NewCollectionHandler() *CollectionHandler {
	return &CollectionHandler{}
}

// Init initializes the handler with database connection
func (h *CollectionHandler) Init() {
	h.service = services.NewCollectionService(database.GetDB())
}

// Create creates a new collection
func (h *CollectionHandler) Create(collection models.Collection) (*models.Collection, error) {
	if err := h.service.Create(&collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

// GetByID retrieves a collection by ID
func (h *CollectionHandler) GetByID(id int64) (*models.Collection, error) {
	return h.service.GetByID(id)
}

// GetAll retrieves all collections
func (h *CollectionHandler) GetAll() ([]models.Collection, error) {
	return h.service.GetAll()
}

// Update updates a collection
func (h *CollectionHandler) Update(collection models.Collection) error {
	return h.service.Update(&collection)
}

// Delete deletes a collection
func (h *CollectionHandler) Delete(id int64) error {
	return h.service.Delete(id)
}

// CreateFolder creates a new folder
func (h *CollectionHandler) CreateFolder(folder models.Folder) (*models.Folder, error) {
	if err := h.service.CreateFolder(&folder); err != nil {
		return nil, err
	}
	return &folder, nil
}

// GetFolderByID retrieves a folder by ID
func (h *CollectionHandler) GetFolderByID(id int64) (*models.Folder, error) {
	return h.service.GetFolderByID(id)
}

// GetFoldersByCollectionID retrieves all folders in a collection
func (h *CollectionHandler) GetFoldersByCollectionID(collectionID int64) ([]models.Folder, error) {
	return h.service.GetFoldersByCollectionID(collectionID)
}

// UpdateFolder updates a folder
func (h *CollectionHandler) UpdateFolder(folder models.Folder) error {
	return h.service.UpdateFolder(&folder)
}

// DeleteFolder deletes a folder
func (h *CollectionHandler) DeleteFolder(id int64) error {
	return h.service.DeleteFolder(id)
}

// GetTree retrieves the full collection tree
func (h *CollectionHandler) GetTree() ([]services.CollectionTree, error) {
	return h.service.GetTree()
}

// MoveRequest moves a request to a different collection/folder
func (h *CollectionHandler) MoveRequest(requestID int64, collectionID int64, folderID *int64) error {
	return h.service.MoveRequest(requestID, collectionID, folderID)
}

// ReorderCollections updates the sort order of collections
func (h *CollectionHandler) ReorderCollections(ids []int64) error {
	return h.service.ReorderCollections(ids)
}

// ReorderFolders updates the sort order of folders in a collection
func (h *CollectionHandler) ReorderFolders(collectionID int64, ids []int64) error {
	return h.service.ReorderFolders(collectionID, ids)
}

// ReorderRequests updates the sort order of requests in a collection/folder
func (h *CollectionHandler) ReorderRequests(collectionID int64, folderID *int64, ids []int64) error {
	return h.service.ReorderRequests(collectionID, folderID, ids)
}
