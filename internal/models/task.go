package models 

import (
	"time"

	"github.com/google/uuid"
)

/*

This is a custom string type with a fixed set of allowed values — 
	an idiom for modeling a finite state machine (const).
*/

type TaskStatus string 

const (
	TaskStatusPending 	TaskStatus = "pending"
	TaskStatusInProcess	TaskStatus = "in_progress"
	TaskStatusPaused	TaskStatus = "paused"
	TaskStatusCanceled	TaskStatus = "canceled"
)

type Task struct {
	ID			uuid.UUID`json:"id" db:"id"`
	Title		string`json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Status		TaskStatus `json:"status" db:"status"`
	ProjectID	uuid.UUID`json:"project_id" db:"project_id"`
	CreatedAt	time.Time`json:"created_at" db:"created_at"`
	UpdatedAt	time.Time`json:"updated_at" db:"updated_at"`
}
