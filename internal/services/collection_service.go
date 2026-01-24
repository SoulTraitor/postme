package services

import (
	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// CollectionService handles collection business logic
type CollectionService struct {
	collectionRepo *repository.CollectionRepository
	folderRepo     *repository.FolderRepository
	requestRepo    *repository.RequestRepository
}

// NewCollectionService creates a new CollectionService
func NewCollectionService(db *sqlx.DB) *CollectionService {
	return &CollectionService{
		collectionRepo: repository.NewCollectionRepository(db),
		folderRepo:     repository.NewFolderRepository(db),
		requestRepo:    repository.NewRequestRepository(db),
	}
}

// Create creates a new collection
func (s *CollectionService) Create(collection *models.Collection) error {
	return s.collectionRepo.Create(collection)
}

// GetByID retrieves a collection by ID
func (s *CollectionService) GetByID(id int64) (*models.Collection, error) {
	return s.collectionRepo.GetByID(id)
}

// GetAll retrieves all collections
func (s *CollectionService) GetAll() ([]models.Collection, error) {
	return s.collectionRepo.GetAll()
}

// Update updates a collection
func (s *CollectionService) Update(collection *models.Collection) error {
	return s.collectionRepo.Update(collection)
}

// Delete deletes a collection
func (s *CollectionService) Delete(id int64) error {
	return s.collectionRepo.Delete(id)
}

// CreateFolder creates a new folder in a collection
func (s *CollectionService) CreateFolder(folder *models.Folder) error {
	return s.folderRepo.Create(folder)
}

// GetFolderByID retrieves a folder by ID
func (s *CollectionService) GetFolderByID(id int64) (*models.Folder, error) {
	return s.folderRepo.GetByID(id)
}

// GetFoldersByCollectionID retrieves all folders in a collection
func (s *CollectionService) GetFoldersByCollectionID(collectionID int64) ([]models.Folder, error) {
	return s.folderRepo.GetByCollectionID(collectionID)
}

// UpdateFolder updates a folder
func (s *CollectionService) UpdateFolder(folder *models.Folder) error {
	return s.folderRepo.Update(folder)
}

// DeleteFolder deletes a folder
func (s *CollectionService) DeleteFolder(id int64) error {
	return s.folderRepo.Delete(id)
}

// MoveFolder moves a folder to a different collection
func (s *CollectionService) MoveFolder(folderID int64, collectionID int64) error {
	return s.folderRepo.Move(folderID, collectionID)
}

// CollectionTree represents a collection with its folders and requests
type CollectionTree struct {
	Collection models.Collection `json:"collection"`
	Folders    []FolderTree      `json:"folders"`
	Requests   []models.Request  `json:"requests"` // Requests directly under collection
}

// FolderTree represents a folder with its requests
type FolderTree struct {
	Folder   models.Folder    `json:"folder"`
	Requests []models.Request `json:"requests"`
}

// GetTree retrieves the full collection tree
func (s *CollectionService) GetTree() ([]CollectionTree, error) {
	collections, err := s.collectionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	requests, err := s.requestRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Build request maps
	requestsByCollection := make(map[int64][]models.Request)
	requestsByFolder := make(map[int64][]models.Request)
	for _, req := range requests {
		if req.FolderID != nil {
			requestsByFolder[*req.FolderID] = append(requestsByFolder[*req.FolderID], req)
		} else {
			requestsByCollection[req.CollectionID] = append(requestsByCollection[req.CollectionID], req)
		}
	}

	var tree []CollectionTree
	for _, col := range collections {
		folders, err := s.folderRepo.GetByCollectionID(col.ID)
		if err != nil {
			return nil, err
		}

		var folderTrees []FolderTree
		for _, folder := range folders {
			folderTrees = append(folderTrees, FolderTree{
				Folder:   folder,
				Requests: requestsByFolder[folder.ID],
			})
		}

		tree = append(tree, CollectionTree{
			Collection: col,
			Folders:    folderTrees,
			Requests:   requestsByCollection[col.ID],
		})
	}

	return tree, nil
}

// MoveRequest moves a request to a different collection/folder
func (s *CollectionService) MoveRequest(requestID int64, collectionID int64, folderID *int64) error {
	request, err := s.requestRepo.GetByID(requestID)
	if err != nil {
		return err
	}
	request.CollectionID = collectionID
	request.FolderID = folderID
	return s.requestRepo.Update(request)
}

// ReorderCollections updates the sort order of collections
func (s *CollectionService) ReorderCollections(ids []int64) error {
	for i, id := range ids {
		collection, err := s.collectionRepo.GetByID(id)
		if err != nil {
			return err
		}
		collection.SortOrder = i
		if err := s.collectionRepo.Update(collection); err != nil {
			return err
		}
	}
	return nil
}

// ReorderFolders updates the sort order of folders in a collection
func (s *CollectionService) ReorderFolders(collectionID int64, ids []int64) error {
	for i, id := range ids {
		folder, err := s.folderRepo.GetByID(id)
		if err != nil {
			return err
		}
		folder.SortOrder = i
		if err := s.folderRepo.Update(folder); err != nil {
			return err
		}
	}
	return nil
}

// ReorderRequests updates the sort order of requests in a collection/folder
func (s *CollectionService) ReorderRequests(collectionID int64, folderID *int64, ids []int64) error {
	for i, id := range ids {
		request, err := s.requestRepo.GetByID(id)
		if err != nil {
			return err
		}
		request.SortOrder = i
		if err := s.requestRepo.Update(request); err != nil {
			return err
		}
	}
	return nil
}
