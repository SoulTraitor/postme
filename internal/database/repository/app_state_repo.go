package repository

import (
	"encoding/json"
	"time"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// AppStateRepository handles app state data access
type AppStateRepository struct {
	db *sqlx.DB
}

// NewAppStateRepository creates a new AppStateRepository
func NewAppStateRepository(db *sqlx.DB) *AppStateRepository {
	return &AppStateRepository{db: db}
}

// Get retrieves the app state
func (r *AppStateRepository) Get() (*models.AppState, error) {
	var state models.AppState
	err := r.db.Get(&state, "SELECT * FROM app_state WHERE id = 1")
	if err != nil {
		return nil, err
	}
	return &state, nil
}

// Update updates the app state
func (r *AppStateRepository) Update(state *models.AppState) error {
	_, err := r.db.Exec(`
		UPDATE app_state SET
			window_width = ?, window_height = ?, window_x = ?, window_y = ?,
			window_maximized = ?, sidebar_open = ?, sidebar_width = ?,
			layout_direction = ?, split_ratio = ?, theme = ?,
			active_env_id = ?, request_timeout = ?, auto_locate_sidebar = ?,
			use_system_proxy = ?, request_panel_tab = ?, updated_at = ?
		WHERE id = 1
	`, state.WindowWidth, state.WindowHeight, state.WindowX, state.WindowY,
		state.WindowMaximized, state.SidebarOpen, state.SidebarWidth,
		state.LayoutDirection, state.SplitRatio, state.Theme,
		state.ActiveEnvID, state.RequestTimeout, state.AutoLocateSidebar,
		state.UseSystemProxy, state.RequestPanelTab, time.Now())
	return err
}

// GetSidebarState retrieves sidebar expanded states
func (r *AppStateRepository) GetSidebarState() ([]models.SidebarState, error) {
	var states []models.SidebarState
	err := r.db.Select(&states, "SELECT * FROM sidebar_state")
	if err != nil {
		return nil, err
	}
	return states, nil
}

// SetSidebarItemExpanded sets the expanded state of a sidebar item
func (r *AppStateRepository) SetSidebarItemExpanded(itemType string, itemID int64, expanded bool) error {
	_, err := r.db.Exec(`
		INSERT INTO sidebar_state (item_type, item_id, expanded) VALUES (?, ?, ?)
		ON CONFLICT(item_type, item_id) DO UPDATE SET expanded = ?
	`, itemType, itemID, expanded, expanded)
	return err
}

// GetTabSessions retrieves all tab sessions
func (r *AppStateRepository) GetTabSessions() ([]models.TabSession, error) {
	var sessions []models.TabSession
	err := r.db.Select(&sessions, "SELECT * FROM tab_sessions ORDER BY sort_order")
	if err != nil {
		return nil, err
	}
	for i := range sessions {
		json.Unmarshal([]byte(sessions[i].HeadersJSON), &sessions[i].Headers)
		json.Unmarshal([]byte(sessions[i].ParamsJSON), &sessions[i].Params)
	}
	return sessions, nil
}

// SaveTabSession saves a tab session
func (r *AppStateRepository) SaveTabSession(session *models.TabSession) error {
	headersJSON, _ := json.Marshal(session.Headers)
	paramsJSON, _ := json.Marshal(session.Params)

	_, err := r.db.Exec(`
		INSERT INTO tab_sessions (tab_id, request_id, title, sort_order, is_active, is_dirty, method, url, headers, params, body, body_type)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(tab_id) DO UPDATE SET
			request_id = ?, title = ?, sort_order = ?, is_active = ?, is_dirty = ?,
			method = ?, url = ?, headers = ?, params = ?, body = ?, body_type = ?, updated_at = CURRENT_TIMESTAMP
	`, session.TabID, session.RequestID, session.Title, session.SortOrder, session.IsActive, session.IsDirty,
		session.Method, session.URL, string(headersJSON), string(paramsJSON), session.Body, session.BodyType,
		session.RequestID, session.Title, session.SortOrder, session.IsActive, session.IsDirty,
		session.Method, session.URL, string(headersJSON), string(paramsJSON), session.Body, session.BodyType)
	return err
}

// DeleteTabSession deletes a tab session
func (r *AppStateRepository) DeleteTabSession(tabID string) error {
	_, err := r.db.Exec("DELETE FROM tab_sessions WHERE tab_id = ?", tabID)
	return err
}

// ClearTabSessions clears all tab sessions
func (r *AppStateRepository) ClearTabSessions() error {
	_, err := r.db.Exec("DELETE FROM tab_sessions")
	return err
}

// SetActiveTab sets the active tab
func (r *AppStateRepository) SetActiveTab(tabID string) error {
	_, err := r.db.Exec("UPDATE tab_sessions SET is_active = 0")
	if err != nil {
		return err
	}
	_, err = r.db.Exec("UPDATE tab_sessions SET is_active = 1 WHERE tab_id = ?", tabID)
	return err
}
