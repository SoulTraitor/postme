package repository

import (
	"encoding/json"
	"time"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// RequestRepository handles request data access
type RequestRepository struct {
	db *sqlx.DB
}

// NewRequestRepository creates a new RequestRepository
func NewRequestRepository(db *sqlx.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

// Create creates a new request
func (r *RequestRepository) Create(req *models.Request) error {
	headersJSON, _ := json.Marshal(req.Headers)
	paramsJSON, _ := json.Marshal(req.Params)

	result, err := r.db.Exec(`
		INSERT INTO requests (collection_id, folder_id, name, method, url, headers, params, body, body_type, sort_order)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.CollectionID, req.FolderID, req.Name, req.Method, req.URL, string(headersJSON), string(paramsJSON), req.Body, req.BodyType, req.SortOrder)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	req.ID = id
	return nil
}

// GetByID retrieves a request by ID
func (r *RequestRepository) GetByID(id int64) (*models.Request, error) {
	var req models.Request
	err := r.db.Get(&req, "SELECT * FROM requests WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	json.Unmarshal([]byte(req.HeadersJSON), &req.Headers)
	json.Unmarshal([]byte(req.ParamsJSON), &req.Params)
	return &req, nil
}

// GetByCollectionID retrieves all requests in a collection
func (r *RequestRepository) GetByCollectionID(collectionID int64) ([]models.Request, error) {
	var requests []models.Request
	err := r.db.Select(&requests, "SELECT * FROM requests WHERE collection_id = ? ORDER BY sort_order", collectionID)
	if err != nil {
		return nil, err
	}

	for i := range requests {
		json.Unmarshal([]byte(requests[i].HeadersJSON), &requests[i].Headers)
		json.Unmarshal([]byte(requests[i].ParamsJSON), &requests[i].Params)
	}
	return requests, nil
}

// GetByFolderID retrieves all requests in a folder
func (r *RequestRepository) GetByFolderID(folderID int64) ([]models.Request, error) {
	var requests []models.Request
	err := r.db.Select(&requests, "SELECT * FROM requests WHERE folder_id = ? ORDER BY sort_order", folderID)
	if err != nil {
		return nil, err
	}

	for i := range requests {
		json.Unmarshal([]byte(requests[i].HeadersJSON), &requests[i].Headers)
		json.Unmarshal([]byte(requests[i].ParamsJSON), &requests[i].Params)
	}
	return requests, nil
}

// Update updates a request
func (r *RequestRepository) Update(req *models.Request) error {
	headersJSON, _ := json.Marshal(req.Headers)
	paramsJSON, _ := json.Marshal(req.Params)

	_, err := r.db.Exec(`
		UPDATE requests SET
			collection_id = ?, folder_id = ?, name = ?, method = ?, url = ?,
			headers = ?, params = ?, body = ?, body_type = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, req.CollectionID, req.FolderID, req.Name, req.Method, req.URL,
		string(headersJSON), string(paramsJSON), req.Body, req.BodyType, req.SortOrder, time.Now(), req.ID)
	return err
}

// Delete deletes a request
func (r *RequestRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}

// GetAll retrieves all requests
func (r *RequestRepository) GetAll() ([]models.Request, error) {
	var requests []models.Request
	err := r.db.Select(&requests, "SELECT * FROM requests ORDER BY collection_id, folder_id, sort_order")
	if err != nil {
		return nil, err
	}

	for i := range requests {
		json.Unmarshal([]byte(requests[i].HeadersJSON), &requests[i].Headers)
		json.Unmarshal([]byte(requests[i].ParamsJSON), &requests[i].Params)
	}
	return requests, nil
}
