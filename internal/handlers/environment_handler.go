package handlers

import (
	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/SoulTraitor/postme/internal/services"
)

// EnvironmentHandler handles environment-related operations for the frontend
type EnvironmentHandler struct {
	service *services.EnvironmentService
}

// NewEnvironmentHandler creates a new EnvironmentHandler
func NewEnvironmentHandler() *EnvironmentHandler {
	return &EnvironmentHandler{}
}

// Init initializes the handler with database connection
func (h *EnvironmentHandler) Init() {
	h.service = services.NewEnvironmentService(database.GetDB())
}

// Create creates a new environment
func (h *EnvironmentHandler) Create(env models.Environment) (*models.Environment, error) {
	if err := h.service.Create(&env); err != nil {
		return nil, err
	}
	return &env, nil
}

// GetByID retrieves an environment by ID
func (h *EnvironmentHandler) GetByID(id int64) (*models.Environment, error) {
	return h.service.GetByID(id)
}

// GetAll retrieves all environments
func (h *EnvironmentHandler) GetAll() ([]models.Environment, error) {
	return h.service.GetAll()
}

// Update updates an environment
func (h *EnvironmentHandler) Update(env models.Environment) error {
	return h.service.Update(&env)
}

// Delete deletes an environment
func (h *EnvironmentHandler) Delete(id int64) error {
	return h.service.Delete(id)
}

// GetGlobalVariables retrieves global variables
func (h *EnvironmentHandler) GetGlobalVariables() (*models.GlobalVariables, error) {
	return h.service.GetGlobalVariables()
}

// UpdateGlobalVariables updates global variables
func (h *EnvironmentHandler) UpdateGlobalVariables(variables []models.Variable) error {
	return h.service.UpdateGlobalVariables(variables)
}
