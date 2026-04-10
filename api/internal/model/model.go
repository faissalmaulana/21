package model

import (
	"time"
)

type Project struct {
	ID   string
	Name string
	// IsArchived is nil if not populated from DB; otherwise true/false.
	IsArchive *bool
	CreatedAt *time.Time
}

type Pagination struct {
	Page             int `json:"page"`
	Size             int `json:"size"`
	TotalItemsInPage int `json:"total_items_in_page"`
	TotalItems       int `json:"total_items"`
	TotalPages       int `json:"total_pages"`
}

type Status int

func (s Status) String() string {
	return []string{"open", "done"}[s]
}

const (
	Open Status = iota
	Done
)

type Task struct {
	ID         string     `json:"id"`
	ProjectID  *string    `json:"project_id"`
	Project    *Project   `json:"project"`
	Name       string     `json:"name"`
	Status     *Status    `json:"status"`
	StartAt    *time.Time `json:"start_at"`
	CreatedAt  *time.Time `json:"created_at"`
	LastUpdate *time.Time `json:"last_update"`
}
