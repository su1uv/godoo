package task

import (
	"time"
	"uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
	DueAt       *time.Time `json:"due_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
