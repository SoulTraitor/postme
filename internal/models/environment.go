package models

import "time"

// Variable represents an environment variable
type Variable struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}

// Environment represents an environment with variables
type Environment struct {
	ID            int64      `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`
	Variables     []Variable `json:"variables" db:"-"`
	VariablesJSON string     `json:"-" db:"variables"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

// GlobalVariables represents global variables (single row)
type GlobalVariables struct {
	ID            int64      `json:"id" db:"id"`
	Variables     []Variable `json:"variables" db:"-"`
	VariablesJSON string     `json:"-" db:"variables"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}
