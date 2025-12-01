package database

import (
	"errors"
)

// SQLDatabase represents a SQL database
type SQLDatabase struct{}

// NewSQLDatabase creates a new SQL database
// TODO: Implement the logic to create a new SQL database
func NewSQLDatabase(config interface{}) (*SQLDatabase, error) {
	if config == nil {
		return nil, errors.New("database: config cannot be nil")
	}
	return &SQLDatabase{}, nil
}
