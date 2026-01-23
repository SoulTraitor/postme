package handlers

import (
	"github.com/SoulTraitor/postme/internal/database"
	"github.com/SoulTraitor/postme/internal/database/repository"
	"github.com/SoulTraitor/postme/internal/models"
)

// AppStateHandler handles app state operations for the frontend
type AppStateHandler struct {
	repo *repository.AppStateRepository
}

// NewAppStateHandler creates a new AppStateHandler
func NewAppStateHandler() *AppStateHandler {
	return &AppStateHandler{}
}

// Init initializes the handler with database connection
func (h *AppStateHandler) Init() {
	h.repo = repository.NewAppStateRepository(database.GetDB())
}

// Get retrieves the app state
func (h *AppStateHandler) Get() (*models.AppState, error) {
	return h.repo.Get()
}

// Update updates the app state
func (h *AppStateHandler) Update(state models.AppState) error {
	return h.repo.Update(&state)
}

// GetSidebarState retrieves sidebar expanded states
func (h *AppStateHandler) GetSidebarState() ([]models.SidebarState, error) {
	return h.repo.GetSidebarState()
}

// SetSidebarItemExpanded sets the expanded state of a sidebar item
func (h *AppStateHandler) SetSidebarItemExpanded(itemType string, itemID int64, expanded bool) error {
	return h.repo.SetSidebarItemExpanded(itemType, itemID, expanded)
}

// GetTabSessions retrieves all tab sessions
func (h *AppStateHandler) GetTabSessions() ([]models.TabSession, error) {
	return h.repo.GetTabSessions()
}

// SaveTabSession saves a tab session
func (h *AppStateHandler) SaveTabSession(session models.TabSession) error {
	return h.repo.SaveTabSession(&session)
}

// DeleteTabSession deletes a tab session
func (h *AppStateHandler) DeleteTabSession(tabID string) error {
	return h.repo.DeleteTabSession(tabID)
}

// ClearTabSessions clears all tab sessions
func (h *AppStateHandler) ClearTabSessions() error {
	return h.repo.ClearTabSessions()
}

// SetActiveTab sets the active tab
func (h *AppStateHandler) SetActiveTab(tabID string) error {
	return h.repo.SetActiveTab(tabID)
}
