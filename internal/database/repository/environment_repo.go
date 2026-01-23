package repository

import (
	"encoding/json"
	"time"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/jmoiron/sqlx"
)

// EnvironmentRepository handles environment data access
type EnvironmentRepository struct {
	db *sqlx.DB
}

// NewEnvironmentRepository creates a new EnvironmentRepository
func NewEnvironmentRepository(db *sqlx.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

// Create creates a new environment
func (r *EnvironmentRepository) Create(env *models.Environment) error {
	variablesJSON, _ := json.Marshal(env.Variables)

	result, err := r.db.Exec(`
		INSERT INTO environments (name, variables)
		VALUES (?, ?)
	`, env.Name, string(variablesJSON))
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	env.ID = id
	return nil
}

// GetByID retrieves an environment by ID
func (r *EnvironmentRepository) GetByID(id int64) (*models.Environment, error) {
	var env models.Environment
	err := r.db.Get(&env, "SELECT * FROM environments WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(env.VariablesJSON), &env.Variables)
	return &env, nil
}

// GetAll retrieves all environments
func (r *EnvironmentRepository) GetAll() ([]models.Environment, error) {
	var envs []models.Environment
	err := r.db.Select(&envs, "SELECT * FROM environments ORDER BY name")
	if err != nil {
		return nil, err
	}
	for i := range envs {
		json.Unmarshal([]byte(envs[i].VariablesJSON), &envs[i].Variables)
	}
	return envs, nil
}

// Update updates an environment
func (r *EnvironmentRepository) Update(env *models.Environment) error {
	variablesJSON, _ := json.Marshal(env.Variables)

	_, err := r.db.Exec(`
		UPDATE environments SET name = ?, variables = ?, updated_at = ?
		WHERE id = ?
	`, env.Name, string(variablesJSON), time.Now(), env.ID)
	return err
}

// Delete deletes an environment
func (r *EnvironmentRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM environments WHERE id = ?", id)
	return err
}

// GetGlobalVariables retrieves global variables
func (r *EnvironmentRepository) GetGlobalVariables() (*models.GlobalVariables, error) {
	var gv models.GlobalVariables
	err := r.db.Get(&gv, "SELECT * FROM global_variables WHERE id = 1")
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(gv.VariablesJSON), &gv.Variables)
	return &gv, nil
}

// UpdateGlobalVariables updates global variables
func (r *EnvironmentRepository) UpdateGlobalVariables(variables []models.Variable) error {
	variablesJSON, _ := json.Marshal(variables)

	_, err := r.db.Exec(`
		UPDATE global_variables SET variables = ?, updated_at = ?
		WHERE id = 1
	`, string(variablesJSON), time.Now())
	return err
}
