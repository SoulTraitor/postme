package repository

import (
	"time"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// FolderRepository handles folder data access
type FolderRepository struct {
	db *sqlx.DB
}

// NewFolderRepository creates a new FolderRepository
func NewFolderRepository(db *sqlx.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

// Create creates a new folder
func (r *FolderRepository) Create(folder *models.Folder) error {
	result, err := r.db.Exec(`
		INSERT INTO folders (collection_id, name, sort_order)
		VALUES (?, ?, ?)
	`, folder.CollectionID, folder.Name, folder.SortOrder)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	folder.ID = id
	return nil
}

// GetByID retrieves a folder by ID
func (r *FolderRepository) GetByID(id int64) (*models.Folder, error) {
	var folder models.Folder
	err := r.db.Get(&folder, "SELECT * FROM folders WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// GetByCollectionID retrieves all folders in a collection
func (r *FolderRepository) GetByCollectionID(collectionID int64) ([]models.Folder, error) {
	var folders []models.Folder
	err := r.db.Select(&folders, "SELECT * FROM folders WHERE collection_id = ? ORDER BY sort_order", collectionID)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

// Update updates a folder
func (r *FolderRepository) Update(folder *models.Folder) error {
	_, err := r.db.Exec(`
		UPDATE folders SET name = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, folder.Name, folder.SortOrder, time.Now(), folder.ID)
	return err
}

// Delete deletes a folder
func (r *FolderRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM folders WHERE id = ?", id)
	return err
}
