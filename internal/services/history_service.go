package services

import (
	"time"

	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// HistoryService handles history business logic
type HistoryService struct {
	repo *repository.HistoryRepository
}

// NewHistoryService creates a new HistoryService
func NewHistoryService(db *sqlx.DB) *HistoryService {
	return &HistoryService{
		repo: repository.NewHistoryRepository(db),
	}
}

// Create creates a new history record
func (s *HistoryService) Create(history *models.History) error {
	history.CreatedAt = time.Now()
	return s.repo.Create(history)
}

// GetAll retrieves all history records
func (s *HistoryService) GetAll() ([]models.History, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a history record by ID
func (s *HistoryService) GetByID(id int64) (*models.History, error) {
	return s.repo.GetByID(id)
}

// Delete deletes a history record
func (s *HistoryService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Clear clears all history records
func (s *HistoryService) Clear() error {
	return s.repo.Clear()
}
