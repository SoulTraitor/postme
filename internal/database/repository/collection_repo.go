package repository

import (
	"time"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// CollectionRepository handles collection data access
type CollectionRepository struct {
	db *sqlx.DB
}

// NewCollectionRepository creates a new CollectionRepository
func NewCollectionRepository(db *sqlx.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

// Create creates a new collection
func (r *CollectionRepository) Create(collection *models.Collection) error {
	result, err := r.db.Exec(`
		INSERT INTO collections (name, description, sort_order)
		VALUES (?, ?, ?)
	`, collection.Name, collection.Description, collection.SortOrder)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	collection.ID = id
	return nil
}

// GetByID retrieves a collection by ID
func (r *CollectionRepository) GetByID(id int64) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.Get(&collection, "SELECT * FROM collections WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// GetAll retrieves all collections
func (r *CollectionRepository) GetAll() ([]models.Collection, error) {
	var collections []models.Collection
	err := r.db.Select(&collections, "SELECT * FROM collections ORDER BY sort_order")
	if err != nil {
		return nil, err
	}
	return collections, nil
}

// Update updates a collection
func (r *CollectionRepository) Update(collection *models.Collection) error {
	_, err := r.db.Exec(`
		UPDATE collections SET name = ?, description = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, collection.Name, collection.Description, collection.SortOrder, time.Now(), collection.ID)
	return err
}

// Delete deletes a collection
func (r *CollectionRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM collections WHERE id = ?", id)
	return err
}
