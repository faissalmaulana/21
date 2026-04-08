package model

import "time"

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
