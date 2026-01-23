package services

import (
	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// RequestService handles request business logic
type RequestService struct {
	repo *repository.RequestRepository
}

// NewRequestService creates a new RequestService
func NewRequestService(db *sqlx.DB) *RequestService {
	return &RequestService{
		repo: repository.NewRequestRepository(db),
	}
}

// Create creates a new request
func (s *RequestService) Create(req *models.Request) error {
	return s.repo.Create(req)
}

// GetByID retrieves a request by ID
func (s *RequestService) GetByID(id int64) (*models.Request, error) {
	return s.repo.GetByID(id)
}

// GetByCollectionID retrieves all requests in a collection
func (s *RequestService) GetByCollectionID(collectionID int64) ([]models.Request, error) {
	return s.repo.GetByCollectionID(collectionID)
}

// GetByFolderID retrieves all requests in a folder
func (s *RequestService) GetByFolderID(folderID int64) ([]models.Request, error) {
	return s.repo.GetByFolderID(folderID)
}

// Update updates a request
func (s *RequestService) Update(req *models.Request) error {
	return s.repo.Update(req)
}

// Delete deletes a request
func (s *RequestService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// GetAll retrieves all requests
func (s *RequestService) GetAll() ([]models.Request, error) {
	return s.repo.GetAll()
}
