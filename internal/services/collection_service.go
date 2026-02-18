package services

import (
	"fmt"
	"time"

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

// ImportCollection creates a new collection from an export file
func (s *CollectionService) ImportCollection(data *models.ExportFile) (*models.Collection, error) {
	// Determine sort order: place at the end
	allCollections, err := s.collectionRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}
	maxSortOrder := 0
	for _, c := range allCollections {
		if c.SortOrder > maxSortOrder {
			maxSortOrder = c.SortOrder
		}
	}

	// Create collection
	collection := &models.Collection{
		Name:        data.Collection.Name,
		Description: data.Collection.Description,
		SortOrder:   maxSortOrder + 1,
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
