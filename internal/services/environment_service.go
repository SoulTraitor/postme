package services

import (
	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// EnvironmentService handles environment business logic
type EnvironmentService struct {
	repo *repository.EnvironmentRepository
}

// NewEnvironmentService creates a new EnvironmentService
func NewEnvironmentService(db *sqlx.DB) *EnvironmentService {
	return &EnvironmentService{
		repo: repository.NewEnvironmentRepository(db),
	}
}

// Create creates a new environment
func (s *EnvironmentService) Create(env *models.Environment) error {
	return s.repo.Create(env)
}

// GetByID retrieves an environment by ID
func (s *EnvironmentService) GetByID(id int64) (*models.Environment, error) {
	return s.repo.GetByID(id)
}

// GetAll retrieves all environments
func (s *EnvironmentService) GetAll() ([]models.Environment, error) {
	return s.repo.GetAll()
}

// Update updates an environment
func (s *EnvironmentService) Update(env *models.Environment) error {
	return s.repo.Update(env)
}

// Delete deletes an environment
func (s *EnvironmentService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// GetGlobalVariables retrieves global variables
func (s *EnvironmentService) GetGlobalVariables() (*models.GlobalVariables, error) {
	return s.repo.GetGlobalVariables()
}

// UpdateGlobalVariables updates global variables
func (s *EnvironmentService) UpdateGlobalVariables(variables []models.Variable) error {
	return s.repo.UpdateGlobalVariables(variables)
}
