package repository

import (
	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// HistoryRepository handles history data access
type HistoryRepository struct {
	db *sqlx.DB
}

// NewHistoryRepository creates a new HistoryRepository
func NewHistoryRepository(db *sqlx.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

// Create creates a new history record
func (r *HistoryRepository) Create(history *models.History) error {
	result, err := r.db.Exec(`
		INSERT INTO history (request_id, method, url, request_headers, request_body, status_code, response_headers, response_body, duration_ms, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, history.RequestID, history.Method, history.URL, history.RequestHeaders, history.RequestBody,
		history.StatusCode, history.ResponseHeaders, history.ResponseBody, history.DurationMs, history.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	history.ID = id

	// Clean up old records if over limit
	r.cleanup()
	return nil
}

// GetAll retrieves all history records
func (r *HistoryRepository) GetAll() ([]models.History, error) {
	var history []models.History
	err := r.db.Select(&history, "SELECT * FROM history ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	return history, nil
}

// GetByID retrieves a history record by ID
func (r *HistoryRepository) GetByID(id int64) (*models.History, error) {
	var history models.History
	err := r.db.Get(&history, "SELECT * FROM history WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// Delete deletes a history record
func (r *HistoryRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM history WHERE id = ?", id)
	return err
}

// Clear clears all history records
func (r *HistoryRepository) Clear() error {
	_, err := r.db.Exec("DELETE FROM history")
	return err
}

// cleanup removes old records if over the limit
func (r *HistoryRepository) cleanup() {
	r.db.Exec(`
		DELETE FROM history WHERE id NOT IN (
			SELECT id FROM history ORDER BY created_at DESC LIMIT ?
		)
	`, models.MaxHistoryRecords)
}
