package dto

import (
	"time"
)

type PostTask struct {
	Name      string    `json:"name" validate:"required,max=255"`
	ProjectID string    `json:"project_id" validate:"required,uuid"`
	StartAt   time.Time `json:"start_at" validate:"required"`
}

type UpdateTask struct {
	ID        string    `param:"id"`
	Name      string    `json:"name" validate:"omitempty,max=255"`
	StartAt   time.Time `json:"start_at" validate:"omitempty"`
	Status    string    `json:"status" validate:"omitempty,oneof=open,done"`
	ProjectID string    `json:"project_id" validate:"omitempty"`
}
