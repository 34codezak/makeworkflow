package models 

import (
	"time"

	// The de facto standard library for generating and parsing UUIDs (128-bit unique identifiers)
	"github.com/google/uuid" // 
)

// Project entity, with struct tags 
// 	for JSON serialization and database mapping.
type Project struct {
	ID			uuid.UUID `json:"id" db:"id"`
	Name		string	`json:"name" db:"name"`
	Description	string 	`json:"description" db:"description"`
	CreatedAt	time.Time	`json:"created_at" db:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at" db:"updated_at"`
}