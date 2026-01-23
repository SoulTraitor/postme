package handlers

import (
	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/SoulTraitor/postme/internal/services"
)

// HistoryHandler handles history-related operations for the frontend
type HistoryHandler struct {
	service *services.HistoryService
}

// NewHistoryHandler creates a new HistoryHandler
func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{}
}

// Init initializes the handler with database connection
func (h *HistoryHandler) Init() {
	h.service = services.NewHistoryService(database.GetDB())
}

// GetAll retrieves all history records
func (h *HistoryHandler) GetAll() ([]models.History, error) {
	return h.service.GetAll()
}

// GetByID retrieves a history record by ID
func (h *HistoryHandler) GetByID(id int64) (*models.History, error) {
	return h.service.GetByID(id)
}

// Delete deletes a history record
func (h *HistoryHandler) Delete(id int64) error {
	return h.service.Delete(id)
}

// Clear clears all history records
func (h *HistoryHandler) Clear() error {
	return h.service.Clear()
}

// Create creates a new history record
func (h *HistoryHandler) Create(history models.History) (*models.History, error) {
	err := h.service.Create(&history)
	if err != nil {
		return nil, err
	}
	return &history, nil
}
