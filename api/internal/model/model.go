package model

import "time"

type Project struct {
	ID   string
	Name string
	// IsArchived is nil if not populated from DB; otherwise true/false.
	IsArchive *bool
	CreatedAt *time.Time
}
